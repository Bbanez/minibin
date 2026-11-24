pub mod math;
pub mod schema;
pub mod utils;

fn main() {
    let args = utils::args::get();
    utils::args::print_help();
    panic!("Invalid arguments combination provided.");
}
