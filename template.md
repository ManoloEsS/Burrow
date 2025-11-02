Create a bubbletea cli app template that has the next characteristics:
- a top stylized  border and the word "Burrow" on the right corner
-Three tabs that can be focused by pressing the keys control-n for next, and control-p for the previous. they must work as a crousel

On the first tab:
A vertical split with the left split containing:
-A horizontal list with four options that can get focused with control-m and switched with h and l as a carousel
-A text box underneath it with room for a title above it, this text box will be focused with control-u
-Another text box underneath it with room for another title above it, this text box will be focused with control-p
-Another text box with with room for a title and a three option horizontal list. The text box will be focused with control-a, and subsequent presses of control-a will change the option on the list
-Another text box with room for a title above it
-Another text box with room for a title above it and an eight option horizontal list. The text box will be focused with control-b and subsequent presses of control-b will change the option on the list
The right side of the split will have:
-Room for a title on top
-two small text displays next to each other with room for a title on their left
-a viewport with room for a title above it.
-another viewport with room for a title above it 
-another viewport with room for a title above it
These two sides can be changed focus with control-r and on the right side the viewports can change focus with tab

leave the other two tabs empty and put placeholder titles. make it so that each text box input is assigned to a variable,
and in the textboxes with horizontal lists above them will assign a variable with a boolean true if selected

when control-s is pressed the values will be assigned to the variables but the values will stay on the boxes
if control space is pressed a command will be run, just leave the command empty for now
