import { mkdtemp, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

type Check = {
    name: string;
    passed: boolean;
};

type Target = {
    name: string;
    checks: Check[];
};

function run(target: Target, name: string, command: string[]): boolean {
    console.debug(`==> ${target.name}: ${name}`);
    const result = Bun.spawnSync(command, {
        stdio: ['inherit', 'inherit', 'inherit'],
    });
    const passed = result.exitCode === 0;
    target.checks.push({ name, passed });
    return passed;
}

function generate(target: Target, language: string): boolean {
    const output = `tests/dist/${language}`;
    return (
        run(target, 'Clear generated schemas', [
            'go',
            'run',
            'main.go',
            '-o',
            output,
            '-clear',
            'true',
        ]) &&
        run(target, 'Generate schemas', [
            'go',
            'run',
            'main.go',
            '-o',
            output,
            '-i',
            'tests/data',
            '-l',
            language,
        ])
    );
}

async function tsFiles(path: string): Promise<string[]> {
    const files: string[] = [];
    for await (const file of new Bun.Glob('*.ts').scan({ cwd: path })) {
        files.push(join(path, file));
    }
    return files;
}

async function testGo(): Promise<Target> {
    const target: Target = { name: 'Go', checks: [] };
    if (generate(target, 'go')) {
        run(target, 'AllTypesDocument round trip', [
            'go',
            'test',
            './tests/go',
            '-run',
            '^TestAllTypesDocumentRoundTrip$',
        ]);
        run(target, 'AllTypesDocumentBatch round trip', [
            'go',
            'test',
            './tests/go',
            '-run',
            '^TestAllTypesDocumentBatchRoundTrip$',
        ]);
        run(target, 'Accessors, mutation, copy, and round trip', [
            'go',
            'test',
            './tests/go',
            '-run',
            '^TestAllTypesDocumentAccessorsMutationAndCopy$',
        ]);
    }
    return target;
}

async function testTypeScript(): Promise<Target> {
    const target: Target = { name: 'TypeScript', checks: [] };
    if (generate(target, 'ts')) {
        run(target, 'AllTypesDocument round trip', [
            'bun',
            'test',
            'tests/ts',
            '--test-name-pattern',
            'round trips every supported type',
        ]);
        run(target, 'AllTypesDocumentBatch round trip', [
            'bun',
            'test',
            'tests/ts',
            '--test-name-pattern',
            'round trips a batch',
        ]);
        run(target, 'Accessors, mutation, copy, and round trip', [
            'bun',
            'test',
            'tests/ts',
            '--test-name-pattern',
            'getters, setters, copy',
        ]);
        run(target, 'Generated code type-check', [
            'bunx',
            'tsc',
            '--noEmit',
            '--target',
            'ESNext',
            '--module',
            'ESNext',
            '--moduleResolution',
            'bundler',
            ...(await tsFiles('tests/dist/ts')),
            ...(await tsFiles('tests/ts')),
        ]);
    }
    return target;
}

async function testCpp(): Promise<Target> {
    const target: Target = { name: 'C++', checks: [] };
    if (!generate(target, 'cpp')) {
        return target;
    }

    const buildDir = await mkdtemp(join(tmpdir(), 'minibin-cpp-'));
    try {
        if (
            run(target, 'Configure CMake', [
                'cmake',
                '-S',
                'tests/cpp',
                '-B',
                buildDir,
            ]) &&
            run(target, 'Build test binary', ['cmake', '--build', buildDir])
        ) {
            run(target, 'AllTypesDocument round trip', [
                join(buildDir, 'minibin_test'),
                'round-trip',
            ]);
            run(target, 'Mutation and property lookup', [
                join(buildDir, 'minibin_test'),
                'mutation',
            ]);
        }
    } finally {
        await rm(buildDir, { recursive: true, force: true });
    }
    return target;
}

function printSummary(targets: Target[]): void {
    const rows = targets.flatMap((target) =>
        target.checks.map((check) => [
            target.name,
            check.name,
            check.passed ? 'PASS' : 'FAIL',
        ]),
    );
    const headers = ['Target', 'Check', 'Result'];
    const widths = headers.map((header, index) =>
        Math.max(header.length, ...rows.map((row) => row[index].length)),
    );
    const divider = `+${widths.map((width) => '-'.repeat(width + 2)).join('+')}+`;
    const line = (values: string[]) =>
        `|${values.map((value, index) => ` ${value.padEnd(widths[index])} `).join('|')}|`;

    console.debug(`\n${divider}`);
    console.debug(line(headers));
    console.debug(divider);
    for (const row of rows) console.debug(line(row));
    console.debug(divider);
}

const targets = await Promise.all([testGo(), testTypeScript(), testCpp()]);
printSummary(targets);

if (targets.some((target) => target.checks.some((check) => !check.passed))) {
    process.exitCode = 1;
} else {
    console.debug('All Go, TypeScript, and C++ checks passed.');
}
