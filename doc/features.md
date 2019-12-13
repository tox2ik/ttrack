# ttrack features

## Clocking

- [x] log timestamps
    - [x] log in, log out
- [x] log custom date
- [x] mark chunk: log out; log in
- [ ] log sick-leave
    - [ ] include sick-leave in calculation
    - [ ] define standard hours
- [x] print duration of last entry when logging out.

## usability

- [x] log to custom file
    
## logic

- [x] wite to file for current month in $TIMETRACK_DIR (~/ttrack)
- [x] auto close/open stamp based on previous stamp
    - [x] in -> out; out -> in
    - [ ] warn if day not same, 
    - [ ] offer to close last day with prompt, auto open for today

## Parsing

- [x] parse relative dates like:
  -[x] -20 minutes
  -[x] yesterday 11:22

## Counting

- [x] calculate average hours 
    - [x] per day
    - [ ] per month
- [ ] add standard-hours for sick-leave days

## Reporting

- [ ] generate a simple (excel) report
- [ ] exclude holidays


## Help

- [ ] document user-flags and commands
    - [ ] in/out mark
    - [ ] custom-time
    
