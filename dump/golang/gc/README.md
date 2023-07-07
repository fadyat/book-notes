### GC

Using `Mark and Sweep` algorithm to implement a GC.

Executes in separate goroutine.

3-Color algorithm using `White`, `Gray`, `Black` to mark the object.

- `White`: Not visited
- `Gray`: Visited but not scanned
- `Black`: Visited and scanned

Deleting all the `White` objects.

- Marking root objects as `Gray`, push them into the `Gray` stack.
- Begins with `Roots` objects (global variables, stack variables) checking for references
to heap objects, pushing new objects into the `Gray` stack.
- Pop the `Gray` stack, mark the object as `Black`.

Main program can change the state of objects, because of that we use `write barrier` to
lock the object when it's being scanned.

Pacer is used to control the GC execution time, when the heap size is larger than the 
2x of the previous heap size, the GC will be triggered.

Pacer also decides which threads will be used for GC, instead of executing app on them.

### Workflow

GC uses ~25% of CPU time, 1/4 of the CPU cores. 
- Configurable with `GOGC` environment variable.

#### Mark setup

- Stop everything, called `stop the world` (STW)
> Waiting until previous GC is finished, all the goroutines are in safe point.
- Enable write barrier

#### Marking 

- Start the world 
- Scanning stack and global variables
> While scanning stack, goroutine is stopped.
- Running 3-color algorithm

#### Mark termination
- Stop the world
- Waiting for all tasks to be finished, cleaning caches, ending the mark 

#### Sweep setup
- Turning off the write barrier
- Start the world
- Cleaning the heap concurrently

#### Sources

- [The GC](https://www.youtube.com/watch?v=gPxFOMuhnUU)
- [Golang GC in general](https://blog.devgenius.io/golang-garbage-collection-in-general-c28ae82558c4)
