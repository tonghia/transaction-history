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
- I am thinking about having 2 versions, the first one focus on making the program work with simple, short data. The second one will focus on handling with a really big file, large data.
- In the first version, I pipelined the process into steps in which the output of one step will be the input of the next step:
  1. Convert the data in csv file into a list of the program data structure.
  2. Filter the list by the period argument.
  3. Calculate total expense and income, and sort the filtered list.
- The first version is easy to test and understand but it also waste lots of memory by transfer data between each step, the default sort function is not efficient too. It cannot work with a large dataset.
- In the second version, I intend to improve 2 things:
  - While reading the data, filter by period argument. When building the filtered list, apply insertion sort.
  - Divide the original file into chunks and process each chunk in parallel. Finally, calculate the total and apply merge sort for each chunk transaction list.

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

There are some attributes which can be suitable for this program:

- Performance: by complete the task fast and can handle a large dataset.
- Testability: ability to easily write test with high coverage.
- Modifiability: ability to change the format of input, ouput of the program.
- Usability: by providing command, flags, or argument to support multiple features.

In the first version, I would choose: Modifiability, Testability.
In the second version, I would choose: Performance and Modifiability.

## Requirement Fullfilment

## Future Work
