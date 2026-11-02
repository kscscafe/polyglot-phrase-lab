use std::env;

fn check(expected: &[String], submitted: &[String]) -> bool {
    expected == submitted
}

fn main() {
    let mut args = env::args().skip(1);
    let expected: Vec<String> = args.next().unwrap_or_default().split('\u{1f}').map(str::to_owned).collect();
    let submitted: Vec<String> = args.next().unwrap_or_default().split('\u{1f}').map(str::to_owned).collect();
    if expected.len() < 2 || submitted.is_empty() {
        eprintln!("arguments contain tokens separated by ASCII unit separator");
        std::process::exit(2);
    }
    println!("{}", if check(&expected, &submitted) { "correct" } else { "incorrect" });
}

#[cfg(test)]
mod tests {
    use super::check;
    #[test]
    fn accepts_exact_order() {
        let words = vec!["我".into(), "喜欢".into(), "中文。".into()];
        assert!(check(&words, &words));
    }
    #[test]
    fn rejects_wrong_order() {
        let expected = vec!["我".into(), "喜欢".into()];
        let submitted = vec!["喜欢".into(), "我".into()];
        assert!(!check(&expected, &submitted));
    }
}
