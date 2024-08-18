mod anki;
mod config;
mod image;
mod util;

use crate::config::Config;
use rusqlite::Connection;
use util::{join_horizontal, Block};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cfg = Config::load()?;

    let image_code = if !cfg.image.path.is_empty() {
        image::path_to_block(cfg.image.path, cfg.image.width, cfg.image.height)?
    } else {
        "".to_owned()
    };

    let conn = Connection::open("/home/licht/.local/share/Anki2/Benutzer 1/collection.anki2")?;

    let card = anki::random_card(conn, 1674145642111)?.render();

    if !image_code.is_empty() {
        let image_block = Block::new(image_code, cfg.image.width, cfg.image.height);
        let card_block = Block::from(card.as_str());

        println!("{}", join_horizontal(vec![&image_block, &card_block]));
    } else {
        println!("{}", card);
    }

    Ok(())
}
