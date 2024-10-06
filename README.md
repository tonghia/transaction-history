# Transaction History

## Thought Process

- Timeline: With 7 days, I could use the first day only to understand the requirement, prepare some ideas. The next 4 days will focus on design, implement, testing, update document cycle. Finally, in the last 2 days, I will review everything again before the submission.
- Initial understanding:
  - The program is for personal use.
  - The main task of this program is process data to get some kind of knowledge.
- Beakdown the assignment:
  - Create a command line application which take 2 arguments: the "period" and a path to a file.
  - Read the file content and convert into our program data structure.
  - Main task: process the data to get the result.
  - Show the result in standard output.
  - Other requirements: have a build system, flexible to change the output format, handle all potential errors approriately.
- Prioritize tasks:
  - The core value of this program is the logic to process the data.
  - The second most value of this program is the program structure and module design.
  - The CLI interface, handling input and output should be completed quickly by using common libraries.
  - The build system, github set up should be completed quickly by using online resources.
- I am thinking about having 2 versions, the first one forcus on making the program work with simple, short data. The second one will focus on handling with a really big file, large data.

## Design Decisions

### Programming Language

This assignment request creating a personal program which will evaluate programming and problem solving skill instead of applying technology or system design, in my opinion, the programming language does not matter. Between 2 preferred languages: Go and Java, Go is more famous for system program and CLI application. Go has better performance and more efficient memory management than Java. In addition, I have better advantage using Go.

### Libraries and Frameworks

In the first version of the program, I go with standard libraries:

- `os`, `flag` for commandline interface.
- `encoding/csv`, `encoding/json`, `path/filepath` for handling input and output files.
- `sort`, `strconv`, `strings`, `time` for handling logic filtering, converting, comparing.
- `testing` package would be use for unit test.

### Quality Attributes

I would choose 2 to 3 from Usability (this is the program for personal usage), Performance (in case I need to process large data), Testability (ensure the quality of the output), and Modifiability (so I could extend the program to help me in more tasks).

## Requirement Fullfilment

## Future Work
