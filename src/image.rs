use std::path::Path;

use base64::{prelude::BASE64_STANDARD, Engine};
use crossterm::{cursor, Command};

fn image_base64<P>(image_path: P) -> String
where
    P: AsRef<Path>,
{
    let image = std::fs::read(image_path).unwrap();
    BASE64_STANDARD.encode(image)
}

pub(crate) fn imagecode_iterm2(
    image_base64: &str,
    width: u32,
    height: u32,
    move_cursor: bool,
) -> String {
    let mut imagecode = String::new();
    imagecode.push_str("\x1b]1337;File=inline=1");

    if width > 0 {
        imagecode.push_str(&format!(";width={}", width));
    }

    if height > 0 {
        imagecode.push_str(&format!(";height={}", height));
    }

    if !move_cursor {
        imagecode.push_str(";doNotMoveCursor=1");
    }

    imagecode.push_str(&format!(":{}\x07", image_base64));

    imagecode
}

pub(crate) fn to_block(
    image_code: &str,
    width: u32,
    height: u32,
    return_cursor: bool,
) -> Result<String, std::fmt::Error> {
    let mut image_block = String::new();

    if return_cursor {
        cursor::SavePosition.write_ansi(&mut image_block)?;
        image_block.push_str(image_code);
        cursor::RestorePosition.write_ansi(&mut image_block)?;
    } else {
        image_block.push_str(image_code);
    }

    let empty_line = format!("{}\n", " ".repeat(width as usize));
    let empty_block = empty_line.repeat(height as usize);

    image_block.push_str(&empty_block);

    Ok(image_block)
}

pub(crate) fn path_to_block<P>(
    image_path: P,
    width: u32,
    height: u32,
) -> Result<String, std::fmt::Error>
where
    P: AsRef<Path>,
{
    let image_base64 = image_base64(image_path);
    let image_code = imagecode_iterm2(&image_base64, width, height, false);

    to_block(&image_code, width, height, false)
}
