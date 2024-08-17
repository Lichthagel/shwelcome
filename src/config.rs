use clap::{Arg, ArgAction, Command};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Default)]
pub(crate) struct ImageConfig {
    pub(crate) path: String,
    pub(crate) width: u32,
    pub(crate) height: u32,
}

#[derive(Serialize, Deserialize, Debug, Default)]
pub(crate) struct Config {
    pub(crate) image: ImageConfig,
}

impl Config {
    pub(crate) fn load() -> Result<Self, Box<dyn std::error::Error>> {
        let mut args: Self = confy::load("shwelcome", "shwelcome")?;

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
            .get_matches();

        if let Some(image_path) = m.get_one::<String>("image_path") {
            args.image.path = image_path.to_owned();
        }

        Ok(args)
    }
}
