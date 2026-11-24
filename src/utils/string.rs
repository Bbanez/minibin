pub fn to_pascal_case(s: &str) -> String {
    let mut s = s[0..1].to_uppercase() + &s[1..];
    let mut i = 0;
    while i < s.len() {
        let mut should_split = false;
        let c = s.chars().nth(i).unwrap();
        if c == '_' || c == '-' || c == ' ' {
            should_split = true;
        }
        if should_split {
            let a = s[0..i].to_string();
            s = a + &s[i + 1..i + 2].to_uppercase() + &s[i + 2..];
            i = i + 2;
        } else {
            i = i + 1;
        }
    }
    s
}
