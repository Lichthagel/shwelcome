use clap::{Arg, ArgAction, Command};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Default)]
pub(crate) struct ImageConfig {
    pub(crate) path: String,
    pub(crate) width: u32,
    pub(crate) height: u32,
}

#[derive(Serialize, Deserialize, Debug, Default)]
pub(crate) struct AnkiConfig {
    pub(crate) db_path: String,
    pub(crate) deck_id: u64,
}

#[derive(Serialize, Deserialize, Debug, Default)]
pub(crate) struct Config {
    pub(crate) image: ImageConfig,
    pub(crate) anki: AnkiConfig,
}

impl Config {
    pub(crate) fn load() -> Result<Self, Box<dyn std::error::Error>> {
        let mut config: Self = confy::load("shwelcome", "shwelcome")?;

        let m = Command::new("shwelcome")
            .author("Lichthagel")
            .version("0.1.0")
            .about("")
            .arg(
                Arg::new("image_path")
                    .short('i')
                    .long("image_path")
                    .value_name("PATH")
                    .action(ArgAction::Set),
            )
            .arg(
                Arg::new("image_width")
                    .long("image_width")
                    .value_name("WIDTH")
                    .action(ArgAction::Set),
            )
            .arg(
                Arg::new("image_height")
                    .long("image_height")
                    .value_name("HEIGHT")
                    .action(ArgAction::Set),
            )
            .arg(
                Arg::new("anki_db_path")
                    .short('a')
                    .long("anki_db_path")
                    .value_name("PATH")
                    .action(ArgAction::Set),
            )
            .arg(
                Arg::new("anki_deck_id")
                    .short('d')
                    .long("anki_deck_id")
                    .value_name("ID")
                    .action(ArgAction::Set),
            )
            .get_matches();

        if let Some(image_path) = m.get_one::<String>("image_path") {
            config.image.path = image_path.to_owned();
        }

        if let Some(image_width) = m.get_one::<u32>("image_width") {
            config.image.width = *image_width;
        }

        if let Some(image_height) = m.get_one::<u32>("image_height") {
            config.image.height = *image_height;
        }

        if let Some(anki_db_path) = m.get_one::<String>("anki_db_path") {
            config.anki.db_path = anki_db_path.to_owned();
        } else if config.anki.db_path.is_empty() {
            eprintln!("Anki database path is required");
            std::process::exit(1);
        }

        if let Some(anki_deck_id) = m.get_one::<u64>("anki_deck_id") {
            config.anki.deck_id = *anki_deck_id;
        } else if config.anki.deck_id == 0 {
            eprintln!("Anki deck ID is required");
            std::process::exit(1);
        }

        Ok(config)
    }
}
