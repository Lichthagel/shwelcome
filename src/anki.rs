use itertools::Itertools;
use rand::seq::SliceRandom;
use rusqlite::Connection;
use yansi::{Color, Paint, Style};

pub(crate) fn random_card(
    connection: Connection,
    deck_id: u64,
) -> Result<AnkiCard, Box<dyn std::error::Error>> {
    let mut stmt = connection.prepare(
        "SELECT notes.flds FROM cards
    JOIN notes ON cards.nid = notes.id
    WHERE cards.did = ?",
    )?;

    let rows = stmt
        .query_map(&[&deck_id], |row| row.get(0))?
        .filter(|r| r.is_ok())
        .map(|r| r.unwrap())
        .collect::<Vec<String>>();

    let card = rows.choose(&mut rand::thread_rng()).unwrap();

    Ok(parse_line(card)?)
}

#[derive(Debug)]
pub(crate) struct AnkiCard {
    word: String,
    reading: String,
    translations: Vec<String>,
}

fn parse_line(line: &str) -> Result<AnkiCard, &'static str> {
    let mut parts = line.split('\x1f');

    let err = "malformed anki card";

    dbg!(line);

    parts.next().ok_or(err)?;

    let word = parts.next().ok_or(err)?.to_string();
    let reading = parts.next().ok_or(err)?.to_string();
    let translations = parts
        .next()
        .ok_or(err)?
        .split("<br>")
        .map(|s| s.trim().to_string())
        .collect();

    Ok(AnkiCard {
        word,
        reading,
        translations,
    })
}

trait RemoveHtmlTags<O> {
    fn remove_html_tags(&self) -> O;
}

impl RemoveHtmlTags<String> for &str {
    fn remove_html_tags(&self) -> String {
        let mut in_tag = false;
        let mut result = String::new();

        for c in self.chars() {
            if c == '<' {
                in_tag = true;
            } else if c == '>' {
                in_tag = false;
            } else if !in_tag {
                result.push(c);
            }
        }

        result
    }
}

impl RemoveHtmlTags<String> for String {
    fn remove_html_tags(&self) -> String {
        self.as_str().remove_html_tags()
    }
}

const fn color_ctp_to_yansi(color: catppuccin::Color) -> Color {
    Color::Rgb(color.rgb.r, color.rgb.g, color.rgb.b)
}

static COLOR_CRUST: Color = color_ctp_to_yansi(catppuccin::PALETTE.mocha.colors.crust);
static COLOR_SUBTEXT0: Color = color_ctp_to_yansi(catppuccin::PALETTE.mocha.colors.subtext0);
static COLOR_PINK: Color = color_ctp_to_yansi(catppuccin::PALETTE.mocha.colors.pink);

fn render_translation(translation: &str) -> String {
    let mut translation = translation.to_string();

    static STYLE_TYPE: Style = Style::new().fg(COLOR_SUBTEXT0).italic();
    static STYLE_NUMBER: Style = Style::new().bg(COLOR_PINK).fg(COLOR_CRUST); // TODO size & alignment
    static STYLE_WORD: Style = Style::new();

    if translation.contains("78909C") {
        translation = translation
            .remove_html_tags()
            .trim()
            .paint(STYLE_TYPE)
            .to_string();
    } else {
        let split: Vec<&str> = translation.splitn(2, "</b>").collect();

        if split.len() == 2 {
            let (number, word) = (split[0], split[1]);

            let number = number.remove_html_tags().trim().to_owned();

            let number = format!("{}{}", " ".repeat(3 - number.len()), number)
                .paint(STYLE_NUMBER)
                .to_string();

            let word = word.remove_html_tags().trim().paint(STYLE_WORD).to_string();

            translation = format!("{} {}", number, word);
        } else {
            translation = translation.paint(STYLE_WORD).to_string();
        }
    }

    return translation;
}

impl AnkiCard {
    pub(crate) fn render(&self) -> String {
        static STYLE_WORD: Style = Style::new().fg(COLOR_PINK);
        static STYLE_READING: Style = Style::new().fg(COLOR_SUBTEXT0);

        let translations = self
            .translations
            .iter()
            .map(|t| render_translation(&t))
            .join("\n");

        format!(
            "{} - {}\n{}",
            self.word.paint(STYLE_WORD),
            self.reading.paint(STYLE_READING),
            translations
        )
    }
}
