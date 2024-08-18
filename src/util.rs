#[derive(Debug, Clone)]
pub(crate) struct Block {
    content: String,
    width: u32,
    height: u32,
}

impl Block {
    pub(crate) fn new(content: String, width: u32, height: u32) -> Self {
        Self {
            content,
            width,
            height,
        }
    }

    fn padded(&self) -> String {
        let actual_height = self.content.lines().count() as u32;

        self.content
            .lines()
            .map(|line| {
                if line.len() < self.width as usize {
                    let padding = self.width as usize - line.len();
                    format!("{}{}", line, " ".repeat(padding))
                } else {
                    line.to_owned()
                }
            })
            .chain(
                std::iter::repeat(" ".repeat(self.width as usize))
                    .take((self.height - actual_height) as usize),
            )
            .collect::<Vec<String>>()
            .join("\n")
    }
}

impl From<&str> for Block {
    fn from(s: &str) -> Self {
        let mut block = Self {
            content: s.to_string(),
            width: 0,
            height: 0,
        };

        for line in s.lines() {
            if line.len() > block.width as usize {
                block.width = line.len() as u32;
            }

            block.height += 1;
        }

        block
    }
}

impl AsRef<Block> for Block {
    fn as_ref(&self) -> &Block {
        self
    }
}

pub(crate) fn join_horizontal<B>(blocks: Vec<B>) -> String
where
    B: AsRef<Block>,
{
    let max_height = blocks.iter().map(|b| b.as_ref().height).max().unwrap_or(0);

    if max_height == 0 {
        return String::new();
    }

    let mut res = String::new();

    for i in 0..max_height {
        for block in &blocks {
            let block = block.as_ref();

            let line = if i < block.height {
                block.padded().lines().nth(i as usize).unwrap().to_owned()
            } else {
                " ".repeat(block.width as usize)
            };

            res.push_str(&line);
        }

        res.push('\n');
    }

    res
}
