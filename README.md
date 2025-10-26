# Date picker bubble

reusable bubble component for date picking

see [pickdate](https://githug.com/maraloon/pickdate) and [tty-diary](https://github.com/maraloo/tty-diary) as examples of using


## TODO

- [x] Show today, style
- [x] Help menu
- [ ] Jumps
    - [x] Jump to today
    - [ ] Month jump
        - [x] p, n
        - [ ] m[1-12]<cr>
    - [ ] Year jump
        - [x] P, N
        - [ ] y[1-12]<cr>
    - [ ] Jump in line: 3l - 3 days later
    - [ ] Jump n month up/down: 3ml/3m<down> - 3 month down
    - [ ] Jump lines: 2j - 2 weeks later
    - [ ] Jump to selected day: `d[1-31]`/`31g`/`31<cr>` will jump on 31th day of current month
- [ ] Lists
    - [ ] Month list (M)
    - [ ] Year list (Y)
- [ ] View
    - [ ] Show 3 month view
    - [ ] Show full year view
    - [ ] Change colors via config
- [ ] Toggle fullsceen (WithAltScreen)
- [x] Center align
- [ ] Toggle week start, monday or sunday
- [ ] CLI opts
    - [x] Week first day
    - [x] Output date format
    - [x] `--start-at date`
    - [ ] Fullscreen
