mod config;
mod image;

use crate::config::Config;
use crossterm::{style::Print, ExecutableCommand};
use std::io::stdout;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cfg = Config::load()?;

    dbg!(&cfg);

    if !cfg.image.path.is_empty() {
        let image_code = image::path_to_block(cfg.image.path, cfg.image.width, cfg.image.height)?;

        stdout().execute(Print(image_code))?;
    }

    Ok(())
}
