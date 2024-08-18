mod anki;
mod config;
mod image;

use crate::config::Config;
use crossterm::{style::Print, ExecutableCommand};
use rusqlite::Connection;
use std::io::stdout;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cfg = Config::load()?;

    dbg!(&cfg);

    if !cfg.image.path.is_empty() {
        let image_code = image::path_to_block(cfg.image.path, cfg.image.width, cfg.image.height)?;

        stdout().execute(Print(image_code))?;
    }

    let conn = Connection::open("/home/licht/.local/share/Anki2/Benutzer 1/collection.anki2")?;

    let card = anki::random_card(conn, 1674145642111)?;

    dbg!(&card);

    println!("{}", card.render());

    Ok(())
}
