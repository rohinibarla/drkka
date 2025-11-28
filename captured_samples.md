First sample, capture.

Replay resulted in wrong output:
```
print "Your table is reserved at the golden Fork for " + PeopleCount + " people on " + Date + " at " + Time + ".Reservation name: " + Name + ", confirmation code: " + Code + "vContact:" + ContactNumber + " 98765-43210. Pl" + TimeDuration + "ase arrive 10 minutes early."
```

```
{
  "examId": "EXAM-DEMO-001",
  "studentId": "3afddf1a-0bdf-468f-ae34-56a811735b36",
  "submissionTime": "2025-11-28T07:54:24.211Z",
  "metadata": {
    "studentName": "Shiva"
  },
  "q1": {
    "questionIndex": 1,
    "questionTitle": "Restaurant Reservation Confirmation",
    "question": "Your table is reserved at The Golden Fork for 2 people on 20-Dec-2023 at 7:30 PM. Reservation name: Rajesh Kumar, Confirmation code: RST892345, Contact: 98765-43210. Please arrive 10 minutes early.",
    "finalAnswer": "print \"Your table is reserved at the golden Fork for \" + PeopleCount + \" people on \" + Date + \" at \" + Time + \".Reservation name: \" + Name + \", confirmation code: \" + Code + \"Contact: \" + ContactNumber + \". Please arrive \" + TimeDuration + \" early.\"",
    "startTime_ms": 1983.8999999761581,
    "endTime_ms": 143767.80000007153,
    "eventLog": [
      {
        "type": "RAW_KEY",
        "key": "p",
        "latency_ms": 0
      },
      {
        "type": "RAW_KEY",
        "key": "r",
        "latency_ms": 129
      },
      {
        "type": "RAW_KEY",
        "key": "i",
        "latency_ms": 75
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 128
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 147
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 225
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 1084
      },
      {
        "type": "RAW_KEY",
        "key": "Y",
        "latency_ms": 508
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 167
      },
      {
        "type": "RAW_KEY",
        "key": "u",
        "latency_ms": 102
      },
      {
        "type": "RAW_KEY",
        "key": "r",
        "latency_ms": 457
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 90
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 442
      },
      {
        "type": "RAW_KEY",
        "key": "a",
        "latency_ms": 141
      },
      {
        "type": "RAW_KEY",
        "key": "b",
        "latency_ms": 84
      },
      {
        "type": "RAW_KEY",
        "key": "l",
        "latency_ms": 160
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 67
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 194
      },
      {
        "type": "RAW_KEY",
        "key": "i",
        "latency_ms": 313
      },
      {
        "type": "RAW_KEY",
        "key": "s",
        "latency_ms": 91
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 84
      },
      {
        "type": "RAW_KEY",
        "key": "r",
        "latency_ms": 142
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 41
      },
      {
        "type": "RAW_KEY",
        "key": "s",
        "latency_ms": 664
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 119
      },
      {
        "type": "RAW_KEY",
        "key": "r",
        "latency_ms": 108
      },
      {
        "type": "RAW_KEY",
        "key": "v",
        "latency_ms": 220
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 155
      },
      {
        "type": "RAW_KEY",
        "key": "d",
        "latency_ms": 143
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 157
      },
      {
        "type": "RAW_KEY",
        "key": "a",
        "latency_ms": 542
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 167
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 91
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 42
      },
      {
        "type": "RAW_KEY",
        "key": "h",
        "latency_ms": 158
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 66
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 100
      },
      {
        "type": "RAW_KEY",
        "key": "g",
        "latency_ms": 376
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 100
      },
      {
        "type": "RAW_KEY",
        "key": "l",
        "latency_ms": 158
      },
      {
        "type": "RAW_KEY",
        "key": "d",
        "latency_ms": 108
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 189
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 145
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 169
      },
      {
        "type": "COMPRESSED",
        "string": "fork",
        "latency_ms": 346,
        "interval_ms": 158
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 520
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 273
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 145
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 122
      },
      {
        "type": "RAW_KEY",
        "key": "F",
        "latency_ms": 529
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 195
      },
      {
        "type": "RAW_KEY",
        "key": "r",
        "latency_ms": 175
      },
      {
        "type": "RAW_KEY",
        "key": "k",
        "latency_ms": 150
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 311
      },
      {
        "type": "RAW_KEY",
        "key": "f",
        "latency_ms": 381
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 66
      },
      {
        "type": "RAW_KEY",
        "key": "r",
        "latency_ms": 117
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 34
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 2153
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 902
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 1089
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 373
      },
      {
        "type": "RAW_KEY",
        "key": "p",
        "latency_ms": 2477
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 399
      },
      {
        "type": "RAW_KEY",
        "key": "P",
        "latency_ms": 379
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 162
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 104
      },
      {
        "type": "RAW_KEY",
        "key": "p",
        "latency_ms": 121
      },
      {
        "type": "RAW_KEY",
        "key": "l",
        "latency_ms": 121
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 106
      },
      {
        "type": "RAW_KEY",
        "key": "C",
        "latency_ms": 435
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 155
      },
      {
        "type": "RAW_KEY",
        "key": "i",
        "latency_ms": 108
      },
      {
        "type": "RAW_KEY",
        "key": "u",
        "latency_ms": 24
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 184
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 376
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 127
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 138
      },
      {
        "type": "RAW_KEY",
        "key": "u",
        "latency_ms": 242
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 175
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 133
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 635
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 313
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 536
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 905
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 267
      },
      {
        "type": "RAW_KEY",
        "key": "p",
        "latency_ms": 229
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 132
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 91
      },
      {
        "type": "RAW_KEY",
        "key": "p",
        "latency_ms": 125
      },
      {
        "type": "RAW_KEY",
        "key": "l",
        "latency_ms": 151
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 107
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 84
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 262
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 369
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 136
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 1028
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 258
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 124
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 1142
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 459
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 231
      },
      {
        "type": "COMPRESSED",
        "string": "Data",
        "latency_ms": 1626,
        "interval_ms": 183
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 937
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 80
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 541
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 531
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 382
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 550
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 264
      },
      {
        "type": "RAW_KEY",
        "key": "a",
        "latency_ms": 193
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 247
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 159
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 1839
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 549
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 398
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 155
      },
      {
        "type": "COMPRESSED",
        "string": "time",
        "latency_ms": 1017,
        "interval_ms": 135
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 421
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 117
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 150
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 135
      },
      {
        "type": "RAW_KEY",
        "key": "T",
        "latency_ms": 424
      },
      {
        "type": "RAW_KEY",
        "key": "i",
        "latency_ms": 245
      },
      {
        "type": "RAW_KEY",
        "key": "m",
        "latency_ms": 146
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 126
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 2881
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 364
      },
      {
        "type": "RAW_KEY",
        "key": ".",
        "latency_ms": 635
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 879
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 685
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 158
      },
      {
        "type": "RAW_KEY",
        "key": ".",
        "latency_ms": 717
      },
      {
        "type": "RAW_KEY",
        "key": "R",
        "latency_ms": 1075
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 158
      },
      {
        "type": "RAW_KEY",
        "key": "s",
        "latency_ms": 575
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 158
      },
      {
        "type": "RAW_KEY",
        "key": "r",
        "latency_ms": 557
      },
      {
        "type": "RAW_KEY",
        "key": "v",
        "latency_ms": 601
      },
      {
        "type": "RAW_KEY",
        "key": "a",
        "latency_ms": 141
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 292
      },
      {
        "type": "RAW_KEY",
        "key": "i",
        "latency_ms": 83
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 100
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 170
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 314
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 183
      },
      {
        "type": "RAW_KEY",
        "key": "a",
        "latency_ms": 100
      },
      {
        "type": "RAW_KEY",
        "key": "m",
        "latency_ms": 100
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 104
      },
      {
        "type": "RAW_KEY",
        "key": ":",
        "latency_ms": 525
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 81
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 825
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 691
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 634
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 340
      },
      {
        "type": "RAW_KEY",
        "key": "N",
        "latency_ms": 1826
      },
      {
        "type": "RAW_KEY",
        "key": "a",
        "latency_ms": 108
      },
      {
        "type": "RAW_KEY",
        "key": "m",
        "latency_ms": 108
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 159
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 579
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 931
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 1096
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 456
      },
      {
        "type": "RAW_KEY",
        "key": ",",
        "latency_ms": 638
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 150
      },
      {
        "type": "RAW_KEY",
        "key": "c",
        "latency_ms": 1001
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 124
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 150
      },
      {
        "type": "RAW_KEY",
        "key": "f",
        "latency_ms": 126
      },
      {
        "type": "RAW_KEY",
        "key": "i",
        "latency_ms": 109
      },
      {
        "type": "RAW_KEY",
        "key": "r",
        "latency_ms": 674
      },
      {
        "type": "RAW_KEY",
        "key": "m",
        "latency_ms": 158
      },
      {
        "type": "RAW_KEY",
        "key": "a",
        "latency_ms": 150
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 184
      },
      {
        "type": "RAW_KEY",
        "key": "i",
        "latency_ms": 75
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 75
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 134
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 171
      },
      {
        "type": "RAW_KEY",
        "key": "c",
        "latency_ms": 877
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 141
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 219
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 602
      },
      {
        "type": "RAW_KEY",
        "key": "d",
        "latency_ms": 123
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 155
      },
      {
        "type": "RAW_KEY",
        "key": ":",
        "latency_ms": 1285
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 61
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 1876
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 530
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 662
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 221
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 522
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 724
      },
      {
        "type": "RAW_KEY",
        "key": "C",
        "latency_ms": 705
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 174
      },
      {
        "type": "RAW_KEY",
        "key": "d",
        "latency_ms": 86
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 142
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 437
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 478
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 364
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 811
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 175,
        "end": 175,
        "latency_ms": 4421
      },
      {
        "type": "RAW_KEY",
        "key": "v",
        "latency_ms": 626
      },
      {
        "type": "RAW_PASTE",
        "content": "Contact: 98765-43210. Please arrive 10 minutes early.",
        "latency_ms": 4
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 1549
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 184,
        "end": 184,
        "latency_ms": 1880
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 619
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 310
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 580
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 330
      },
      {
        "type": "RAW_KEY",
        "key": "C",
        "latency_ms": 581
      },
      {
        "type": "RAW_KEY",
        "key": "o",
        "latency_ms": 232
      },
      {
        "type": "RAW_KEY",
        "key": "n",
        "latency_ms": 117
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 141
      },
      {
        "type": "RAW_KEY",
        "key": "a",
        "latency_ms": 175
      },
      {
        "type": "RAW_KEY",
        "key": "c",
        "latency_ms": 96
      },
      {
        "type": "RAW_KEY",
        "key": "t",
        "latency_ms": 195
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 927
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 1997
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 1149
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 198,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 146
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 199,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 156
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 200,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 160
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 201,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 156
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 202,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 155
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 203,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 164
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 204,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 159
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 205,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 156
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 206,
        "latency_ms": 0
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 155
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 207,
        "latency_ms": 1
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowRight",
        "latency_ms": 141
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 197,
        "end": 208,
        "latency_ms": 1
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 501
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 713
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 195,
        "end": 195,
        "latency_ms": 1519
      },
      {
        "type": "RAW_KEY",
        "key": "N",
        "latency_ms": 703
      },
      {
        "type": "RAW_KEY",
        "key": "u",
        "latency_ms": 219
      },
      {
        "type": "RAW_KEY",
        "key": "m",
        "latency_ms": 211
      },
      {
        "type": "COMPRESSED",
        "string": "ber",
        "latency_ms": 203,
        "interval_ms": 81
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 221,
        "end": 223,
        "latency_ms": 4311
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 1863
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 706
      },
      {
        "type": "RAW_KEY",
        "key": "+",
        "latency_ms": 184
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 225
      },
      {
        "type": "COMPRESSED",
        "string": "+ +",
        "latency_ms": 3345,
        "interval_ms": 173
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 699
      },
      {
        "type": "RAW_KEY",
        "key": "\"",
        "latency_ms": 626
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowLeft",
        "latency_ms": 375
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowLeft",
        "latency_ms": 751
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowLeft",
        "latency_ms": 149
      },
      {
        "type": "RAW_SPECIAL",
        "key": "ArrowLeft",
        "latency_ms": 138
      },
      {
        "type": "RAW_KEY",
        "key": " ",
        "latency_ms": 529
      },
      {
        "type": "COMPRESSED",
        "string": "time",
        "latency_ms": 5673,
        "interval_ms": 147
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 986
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 134
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 172
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 150
      },
      {
        "type": "RAW_KEY",
        "key": "T",
        "latency_ms": 2294
      },
      {
        "type": "RAW_KEY",
        "key": "i",
        "latency_ms": 256
      },
      {
        "type": "RAW_KEY",
        "key": "m",
        "latency_ms": 1014
      },
      {
        "type": "RAW_KEY",
        "key": "e",
        "latency_ms": 160
      },
      {
        "type": "RAW_KEY",
        "key": "D",
        "latency_ms": 959
      },
      {
        "type": "COMPRESSED",
        "string": "uration",
        "latency_ms": 192,
        "interval_ms": 127
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 245,
        "end": 245,
        "latency_ms": 1513
      },
      {
        "type": "SELECTION_CHANGE",
        "start": 242,
        "end": 249,
        "latency_ms": 143
      },
      {
        "type": "RAW_SPECIAL",
        "key": "Backspace",
        "latency_ms": 475
      }
    ]
  }
}
 
```