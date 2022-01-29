# Pomodoro CLI
A simple command line Pomodoro timer.

> The **Pomodoro Technique** is a time management method developed by Francesco Cirillo in the late 1980s. The technique uses a timer to break down work into intervals, traditionally 25 minutes in length, separated by short breaks. Each interval is known as a _pomodoro_, from the Italian word for 'tomato', after the tomato-shaped kitchen timer that Cirillo used as a university student.

### Usage

1. Clone the repository:
  `git clone https://github.com/ihsavru/pomodoro.git`

2. Run the program:
  ```
  cd demo
  go run main.go
  ```

3. You can pass the following flags to it:
  ```
  -longBreak int
    Duration of long break interval in minutes (default 5)
  -shortBreak int
    Duration of short break interval in minutes (default 5)
  -work int
    Duration of work interval in minutes (default 25)
  ```

  **Example:**

  ```
  go run main.go --work=45 --shortBreak=15 --longBreak=45
  ```

5. Press *Ctrl + C* or *q* to quit.

![enter image description here](https://user-images.githubusercontent.com/22816171/84430319-257fa800-ac47-11ea-8c3c-cf808595b2f6.png)
The notification sound is downloaded from freesound.org uploaded by user [RSilveira_88](https://freesound.org/people/RSilveira_88/sounds/216306/).
