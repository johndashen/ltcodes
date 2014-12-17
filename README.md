README
=======

Design:
=======
Generating a random number is abstracted in a RandGen interfaces, of which LinearGen implements the simple linear generator.  The Soliton class allows a random number to be drawn from the robust soliton distribution, by storing the cdf and performing a binary search on the cdf for a given random double in [0,1]. 

LT encoding is handled by BlockEncoder. The file is first stored in a large array of bytes.  An accessor function allows a particular slice to be retrieved over that array.  For constructing the xor of a set of blocks, the first block is copied and then xors are done in place over that block.  This is done in a single thread.

LT decoding is handled by BlockDecoder.  A map of block numbers to confirmed blocks is stored, as well as a list of dirty blocks which are mixed, with the indices of each block that are being x'ored. 	When a new block is added, it is either added to the list of confirmed blocks if it is single, or a dirty block.  When dirty blocks are made clean, they are moved from the list of dirty blocks to the map of clean blocks, and when enough clean blocks are made, the decoder can signal that it is done.  This is done in a single-threaded fashion.  

To my knowledge, LT encoding and decoding are both correct, but LT decoding may require more blocks than normal.
 
Performance:
===========
In practice, LT encoding is performant, but the LT decoding is much slower than the reference implementation because the representation of encoded blocks is straightforward.  The decoder can definitely be improved by changing the index representation of the indices of dirty blocks (which is O(# blocks * mean(soliton distribution))) to a pointer representation to blocks in memory, where each target block maps to a posting list, and each encoded block has a pointer to it from that posting list.  This would be a relatively simple change.

The signal for LT decoding to finish is also dependent upon a user method call, which can be changed by installing a callback in BlockDecoder when the decoder is done. 

Otherwise, it would be possible to speed performance by streaming blocks into a worker pool, and having each worker hold locks only for the target decoding posting lists.  