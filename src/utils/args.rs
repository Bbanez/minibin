#[derive(Debug)]
pub struct Args {
    pub string_templates: bool,
    pub setup: bool,
    pub input: Vec<String>,
    pub output: Vec<String>,
    pub lang: String,
    pub package_base: String,
    pub inject_bson: bool,
    pub clear: bool,
    pub project_build: bool,
    pub create_release: bool,
}

pub struct ArgInfo {
    pub flags: Vec<&'static str>,
    pub typ: &'static str,
    pub desc: &'static str,
    pub allowed_values: Option<Vec<&'static str>>,
    pub parser: fn(argv: &Vec<String>, at_idx: usize, info: &ArgInfo, args: &mut Args),
}

pub fn get_info() -> Vec<ArgInfo> {
    return vec![
        ArgInfo {
            flags: vec!["--input", "-input", "input", "--i", "-i", "i"],
            desc: "Where to find schema files, Default: minibin-schemas",
            typ: "string",
            allowed_values: None,
            parser: |argv: &Vec<String>, at_idx: usize, _info: &ArgInfo, args: &mut Args| {
                args.input = argv[at_idx + 1].split("/").map(|s| s.to_string()).collect();
            },
        },
        ArgInfo {
            flags: vec!["--output", "-output", "output", "--o", "-o", "o"],
            desc: "Where to output generated files",
            typ: "string",
            allowed_values: None,
            parser: |argv: &Vec<String>, at_idx: usize, _info: &ArgInfo, args: &mut Args| {
                args.output = argv[at_idx + 1].split("/").map(|s| s.to_string()).collect();
            },
        },
        ArgInfo {
            flags: vec![
                "--language",
                "-language",
                "language",
                "--lang",
                "-lang",
                "lang",
                "--l",
                "-l",
                "l",
            ],
            desc: "For which language to generate code, Default: go",
            typ: "string",
            allowed_values: Some(vec!["go", "rust", "cpp", "ts"]),
            parser: |argv: &Vec<String>, at_idx: usize, info: &ArgInfo, args: &mut Args| {
                args.lang = argv[at_idx + 1].to_string();
                if !info
                    .allowed_values
                    .as_ref()
                    .unwrap()
                    .contains(&args.lang.as_str())
                {
                    panic!(
                        "Invalid value for {}: {}, allowed values: {:?}",
                        argv[at_idx], args.lang, info.allowed_values
                    );
                }
            },
        },
        ArgInfo {
            flags: vec!["--package", "-package", "package", "--pkg", "-pkg", "pkg"],
            desc: "Base package name for generated code, Default: empty",
            typ: "string",
            allowed_values: None,
            parser: |argv: &Vec<String>, at_idx: usize, _info: &ArgInfo, args: &mut Args| {
                args.package_base = argv[at_idx + 1].to_string();
            },
        },
        ArgInfo {
            flags: vec!["--bson", "-bson", "bson"],
            desc: "Base package name for generated code, Default: empty",
            typ: "bool",
            allowed_values: None,
            parser: |_argv: &Vec<String>, _at_idx: usize, _info: &ArgInfo, args: &mut Args| {
                args.inject_bson = true;
            },
        },
        ArgInfo {
            flags: vec!["--clear", "-clear", "clear", "--c", "-c", "c"],
            desc: "Clear minibin output directory",
            typ: "bool",
            allowed_values: None,
            parser: |_argv: &Vec<String>, _at_idx: usize, _info: &ArgInfo, args: &mut Args| {
                args.clear = true;
            },
        },
        ArgInfo {
            flags: vec![
                "--string-templates",
                "-string-templates",
                "string-templates",
                "--st",
                "-st",
                "st",
            ],
            desc: "INTERNAL: Generate string templates",
            typ: "bool",
            allowed_values: None,
            parser: |_argv: &Vec<String>, _at_idx: usize, _info: &ArgInfo, args: &mut Args| {
                args.string_templates = true;
            },
        },
    ];
}

pub fn get() -> Args {
    let argv: Vec<String> = std::env::args().collect();
    let mut args = Args {
        string_templates: false,
        setup: false,
        input: vec!["minibin-schemas".to_string()],
        output: vec!["minibin-out".to_string()],
        lang: "".to_string(),
        package_base: "".to_string(),
        inject_bson: false,
        clear: false,
        project_build: false,
        create_release: false,
    };
    let arg_infos = get_info();
    for i in 0..arg_infos.len() {
        for idx in 0..argv.len() {
            let arg_value = &argv[idx];
            if !arg_infos[i].flags.contains(&arg_value.as_str()) {
                continue;
            }
            (arg_infos[i].parser)(&argv, idx, &arg_infos[i], &mut args);
        }
    }
    return args;
}

pub fn print_help() {
    let mut message = "".to_string()
        + "┌---------------------------------------------------------------┐\n"
        + "│                             HELP                              │\n"
        + "└---------------------------------------------------------------┘\n\n";
    let arg_infos = get_info();
    for i in 0..arg_infos.len() {
        let info = &arg_infos[i];
        if info.desc.contains("INTERNAL") {
            continue;
        }
        message = message + "---- " + info.desc + "\n";
        message = message + "        Type: " + info.typ + "\n";
        message = message + "        Flags: " + info.flags.join(" │ ").as_str() + "\n\n";
    }
    println!("{}", message);
}
