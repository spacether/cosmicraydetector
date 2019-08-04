# cosmicraydetector
A Go program to detect bit flips caused by cosmic rays

## Implementation Notes
This program makes an array of unsigned 64 bit integers taking up a certain amount of memory (1 GiB default). Every delay seconds (120 default) we inspect all values in the array, and if any have changed we store the bit flip information and display it in the program.

An example of the bit flip information that we store is:
```
f := flip{
	Value: 2,
	Binary: "00000010",
	NumChangedBits: 1,
	ChangedBits: "______X_"
	Duration: 2.307184703s  // how long the value was stored before it was changed
	Time: 2019-08-02 13:29:15.719184 -0700 PDT m=+2.409464913  // when the bit flip happened
}
```


## Usage
`make run` to run the program

## Testing
`make test` to run tests