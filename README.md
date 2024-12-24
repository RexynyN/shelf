# Shelf 
*A nifty CLI tool for the file system power user* 

--- 

### Status

***There are no stable releases at the moment, use the source code at your own risk.***

### TODOS: 

- Polish the rename function
    - Fix the rename operations saving and reverting
    - Fix the duplicate name issue before commiting the filename changes
- Implement the "duplicate" function
    - Search for a consistent and reliable way to find a duplicate - **DONE**
    - Search for a consistent algorithm to find partial duplicates - 
    - Implement the subcommand "dir" to check duplicates between directories  
- Implement Named Duplicates
    - HashMap -> Number of times that name has apperead
        - Dupped-dup name ("image (1).jpg" and "image.jpg" must be seen as equal)
        - Normalize name (optional) -> "IMAGE.jpg" and "image.jpg" must be the same 
- Implement the "tidy" function
    - Search for a consistent algorithm to group a bunch of files in useful folders
        - By Name
        - By Extension
        - By Semantic (Photos, Executables, Videos, Text)
- Implement the "clog" function
    - Find big files that clog up memory in a folder recursively
- Find out even more cool functions for the cool user B^)

### 