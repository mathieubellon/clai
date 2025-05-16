use crossterm::{
    cursor,
    event::{self, Event, KeyCode, KeyModifiers},
    execute,
    style::{self, Stylize},
    terminal::{self, ClearType},
};
use std::io::stdout;
use std::io::{self, Write};
use std::process::Command;

fn main() {
    // Get the sentence from command line arguments
    let args: Vec<String> = std::env::args().collect();
    let sentence = if args.len() > 1 {
        args[1].clone()
    } else {
        println!("Please enter a sentence:");
        let mut input = String::new();
        io::stdin()
            .read_line(&mut input)
            .expect("Failed to read line");
        input.trim().to_string()
    };

    // Generate 3 options based on the sentence
    let options = vec![
        format!("Option 1: {}", "ls -la"),
        format!("Option 2: {}", sentence.to_lowercase()),
        format!("Option 3: {}", sentence.chars().rev().collect::<String>()),
    ];

    // Display options and handle selection
    match select_option(&options) {
        Some(selected) => {
            // Clean up terminal
            let _stdout = stdout();
            terminal::disable_raw_mode().unwrap();

            // Extract the actual content after the "Option X: " prefix
            let content = selected
                .split(": ")
                .collect::<Vec<&str>>()
                .get(1)
                .unwrap_or(&selected.as_str())
                .to_string();

            println!("Selected: {}", content);
            println!("Press Enter to execute this as a command...");

            // Wait for Enter key
            let mut input = String::new();
            io::stdin()
                .read_line(&mut input)
                .expect("Failed to read line");

            // Execute the selected content as a command
            println!("Executing: {}", content);
            match Command::new("sh").arg("-c").arg(&content).status() {
                Ok(status) => println!("Command executed with status: {}", status),
                Err(e) => eprintln!("Failed to execute command: {}", e),
            }
        }
        None => eprintln!("Selection cancelled"),
    }
}

fn select_option(options: &[String]) -> Option<String> {
    let mut stdout = stdout();

    // Enter raw mode
    terminal::enable_raw_mode().unwrap();
    execute!(stdout, terminal::EnterAlternateScreen, cursor::Hide).unwrap();

    let mut selected = 0;

    // Display initial options
    render_options(options, selected);

    // Handle key input
    loop {
        if let Event::Key(key) = event::read().unwrap() {
            match key.code {
                KeyCode::Up => {
                    if selected > 0 {
                        selected -= 1;
                        render_options(options, selected);
                    }
                }
                KeyCode::Down => {
                    if selected < options.len() - 1 {
                        selected += 1;
                        render_options(options, selected);
                    }
                }
                KeyCode::Enter => {
                    // Clean up and exit raw mode
                    execute!(
                        stdout,
                        terminal::Clear(ClearType::All),
                        cursor::Show,
                        terminal::LeaveAlternateScreen
                    )
                    .unwrap();
                    terminal::disable_raw_mode().unwrap();
                    return Some(options[selected].clone());
                }
                KeyCode::Char('c') if key.modifiers.contains(KeyModifiers::CONTROL) => {
                    // Clean up and exit raw mode
                    execute!(
                        stdout,
                        terminal::Clear(ClearType::All),
                        cursor::Show,
                        terminal::LeaveAlternateScreen
                    )
                    .unwrap();
                    terminal::disable_raw_mode().unwrap();
                    return None;
                }
                _ => {}
            }
        }
    }
}

fn render_options(options: &[String], selected: usize) {
    let mut stdout = stdout();

    // Clear the screen
    execute!(
        stdout,
        terminal::Clear(ClearType::All),
        cursor::MoveTo(0, 0)
    )
    .unwrap();

    // Display options
    for (i, option) in options.iter().enumerate() {
        if i == selected {
            execute!(stdout, style::PrintStyledContent(option.clone().bold())).unwrap();
        } else {
            execute!(stdout, style::Print(option)).unwrap();
        }
        execute!(stdout, style::Print("\n")).unwrap();
    }

    execute!(
        stdout,
        style::Print("\nUse up/down arrows to navigate, Enter to select\n")
    )
    .unwrap();
    stdout.flush().unwrap();
}
