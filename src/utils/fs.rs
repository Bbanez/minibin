pub struct FileTreeNode {
    pub abs_path: String,
    pub rel_path: String,
    pub dir_path: String,
}

pub struct FS {
    pub is_win: bool,
    pub slash: String,
    pub base_path: String,
    pub base_len: i32,
}

impl FS {
    pub fn new(base_path: Option<Vec<String>>) -> FS {
        let mut is_win = false;
        let mut slash = "/".to_string();
        if cfg!(target_os = "windows") {
            is_win = true;
            slash = "\\".to_string();
        }
        let mut bp = cwd();
        if let Some(bp_vec) = base_path {
            bp = bp_vec.join(&slash);
        }
        let base_len: Vec<&str> = bp.split(&slash).collect();
        let fs = FS {
            is_win,
            slash,
            base_path: bp.to_string(),
            base_len: base_len.len() as i32,
        };
        fs
    }

    pub fn join(&self, path: &[&str]) -> String {
        path.join(&self.slash)
    }

    pub fn path(&self, path: &[&str]) -> String {
        let p = self.join(path);
        if p.starts_with("/") || p.contains("\\\\") {
            return p;
        }
        return self.base_path.as_str().to_owned() + self.slash.as_str() + &p;
    }

    pub fn to_path(&self, path: &str) -> Vec<String> {
        let parts: Vec<&str> = path.split(&self.slash).collect();
        let mut relative_path: Vec<String> = vec![];
        for i in 0..parts.len() {
            let part = parts[i];
            if part == "" {
                continue;
            }
            relative_path.push(part.to_string());
        }
        return relative_path;
    }

    pub fn file_size(&self, path: &[&str]) -> Result<i64, std::io::Error> {
        let p = path.join(&self.slash);
        let file_info = std::fs::metadata(p)?;
        return Ok(file_info.len() as i64);
    }

    pub fn exists(&self, path: &[&str]) -> bool {
        let p = path.join(&self.slash);
        return std::path::Path::new(&p).exists();
    }

    pub fn mkdir_if_not_exists(&self, path: &str, skip_root: bool) {
        let parts: Vec<&str> = path.split(&self.slash).collect();
        let mut start_at: usize = 0;
        if skip_root {
            start_at = self.base_len as usize;
        }
        for i in start_at..parts.len() {
            let part = parts[i];
            if part == "" {
                continue;
            }
            let path = parts[0..=i].join(&self.slash);
            if !dir_exists(&path) {
                std::fs::create_dir(&path).unwrap();
            }
        }
    }

    pub fn mkdir(&self, path: &[&str]) -> Result<(), std::io::Error> {
        let a_path = self.path(path);
        if dir_exists(&a_path) {
            return Ok(());
        }
        std::fs::create_dir_all(&a_path)?;
        return Ok(());
    }

    pub fn write(&self, data: &[u8], path: &[&str]) -> Result<(), std::io::Error> {
        let a_path = self.path(path);
        self.mkdir_if_not_exists(&a_path, true);
        std::fs::write(&a_path, data)?;
        return Ok(());
    }

    pub fn append(&self, data: &[u8], path: &[&str]) -> Result<bool, std::io::Error> {
        let a_path = self.path(path);
        self.mkdir_if_not_exists(&a_path, true);
        let mut file = std::fs::OpenOptions::new();
        file.append(true).open(&a_path)?;
        std::fs::write(&a_path, data)?;
        return Ok(true);
    }

    pub fn read(&self, path: &[&str]) -> Result<Vec<u8>, std::io::Error> {
        let a_path = self.path(path);
        let data = std::fs::read(&a_path)?;
        return Ok(data);
    }

    pub fn read_string(&self, path: &[&str]) -> Result<String, std::io::Error> {
        let a_path = self.path(path);
        let data = std::fs::read_to_string(&a_path)?;
        return Ok(data);
    }

    pub fn open_file(&self, path: &[&str]) -> Result<std::fs::File, std::io::Error> {
        let a_path = self.path(path);
        let file = std::fs::File::open(&a_path)?;
        return Ok(file);
    }

    pub fn delete(&self, path: &[&str]) -> Result<bool, std::io::Error> {
        let a_path = self.path(path);
        if !self.exists(path) {
            return Ok(false);
        }
        std::fs::remove_file(&a_path)?;
        return Ok(true);
    }

    pub fn file_tree(
        &self,
        start_location: &[&str],
        current_location: &[&str],
    ) -> Result<Vec<FileTreeNode>, std::io::Error> {
        let start_loc = self.path(start_location);
        let curr_loc: String;
        if current_location.len() == 0 {
            curr_loc = "".to_string();
        } else {
            curr_loc = self.join(current_location) + self.slash.as_str();
        }
        let mut output: Vec<FileTreeNode> = vec![];
        let files = std::fs::read_dir(&start_loc)?;
        for file in files {
            let f = file?;
            let file_info = f.metadata()?;
            let filename = f.file_name().into_string().unwrap();
            let file_path =
                start_loc.clone() + self.slash.as_str() + curr_loc.as_str() + filename.as_str();
            if file_info.is_dir() {
                let children = self.file_tree(
                    start_location,
                    &[current_location, &[filename.as_str()]].concat(),
                );
                match children {
                    Ok(children) => {
                        output.extend(children);
                    }
                    Err(e) => {
                        return Err(e);
                    }
                }
            } else {
                let rel_path = file_path.replace(start_loc.as_str(), "");
                let rel_path = rel_path.trim_prefix("/");
                output.push(FileTreeNode {
                    abs_path: file_path,
                    rel_path: rel_path.to_string(),
                    dir_path: curr_loc.clone(),
                });
            }
        }
        Ok(output)
    }
}

pub fn cwd() -> String {
    return std::env::current_dir()
        .unwrap()
        .to_str()
        .unwrap()
        .to_string();
}

fn dir_exists(path: &str) -> bool {
    let info = std::fs::metadata(path);
    match info {
        Ok(info) => info.is_dir(),
        Err(_) => false,
    }
}
