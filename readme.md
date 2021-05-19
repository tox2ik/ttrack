# ttrack

A cli utility for tracking time.

## Todo


- [ ] dry run with -n => show what the state would be after applying the current stamp
- [ ] write tests
    - [ ] ttrack '07:00'; ttrack '18:00'
    - [ ] ttrack 'yesterday 09:00'; ttrack 'yesterday 17:00';
    - [ ] mark from 'out' state; seems like it adds in but not out.
- [ ] e for edit (stamp file); el edit-log (log file)
- [ ] ask for confirmation when writing to a new file that does not look like a month
- [ ] t or today for listing current stamps
- [ ] parse "simple" fuzzy dates with a few formats
      https://stackoverflow.com/questions/41329583/relative-date-parsing
      https://github.com/tj/go-naturaldate/blob/master/naturaldate_test.ghttps://cto.ai/static/hero-visual-centered-slack.pngo
- [ ] continue verb to add a samp one-second after the last
- [ ] stamp / log sick-leave
- [ ] improve '20 min ago' behavior -> imply today
- [ ] report error on invalid date e.g. tt 'an hour ago'
- [ ] test tt s output
- [ ] test tt c output
- [ ] tui: action for splitting over-midnight stamps
- [ ] tui: edit log entries


## Usage


Synopsis:

`tt [in|out] [date-time specification]`
`tt [mark|count|sum|log]`

All parameters are optional.


### Typical log in / log out flow:

    $ tt
    in:  2019-12-13 19:28:38 1576261718 -> /home/jaroslav/ttrack/dec

    $ tt
    out: 2019-12-13 20:28:56 1576261736 -> /home/jaroslav/ttrack/dec
    2019-12-13  1.00

### Mark

    $ tt mark

Log out and back in (to mark some time as spent on a specific task).


### Log

Append the hours difference from the last completed entry into a log file.
You may edit the log file and describe what you did.

    $ bash src/tt/tt.sh log
    2019-12-13  0.00: describe activity...

    $ cat ~/ttrack/dec.log
    2019-12-13  0.00: describe activity...


### Count / per day

    $ tt count
    2019-12-13  0.00
    2019-12-13  1.06
    2019-12-13  0.83
        total:  1.90
      average:  0.63

    $ tt count per-day
    2019-12-13  1.90
        total:  1.90
      average:  1.90

### Custom time

Time parsing in handled by `date -d` from GNU Coreutils unless specified as `yyyy-mm-dd hh:mm:ss`.

## Typical records file


    $ cat ~/ttrack/nov | sed -e s/inn/in/ -e s/ut/out/
    in:  2019-11-04 09:50:00 1572857400
    out: 2019-11-04 12:23:01 1572866581
    in:  2019-11-04 12:23:01 1572866581
    out: 2019-11-04 18:29:57 1572888597
    in:  2019-11-05 09:45:00 1572943500
    out: 2019-11-05 11:27:12 1572949632
    in:  2019-11-05 11:27:12 1572949632
    out: 2019-11-05 17:44:34 1572972274
    in:  2019-11-06 08:22:55 1573024975
    out: 2019-11-06 09:45:10 1573029910
    in:  2019-11-06 09:45:10 1573029910
    out: 2019-11-06 11:25:12 1573035912
    in:  2019-11-06 11:25:12 1573035912
    out: 2019-11-06 13:32:57 1573043577
    in:  2019-11-06 13:32:57 1573043577
    out: 2019-11-06 15:51:27 1573051887
    in:  2019-11-06 15:51:27 1573051887
    out: 2019-11-06 16:12:03 1573053123
    in:  2019-11-06 16:12:04 1573053124
    out: 2019-11-06 18:32:37 1573061557

## Console Ui toolkits

- https://appliedgo.net/tui/
- https://github.com/jroimartin/gocui
