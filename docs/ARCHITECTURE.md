// DATA - SQL //

Things can be changed through website/admin listed here, everything else mentioned in frontend is content in code instead of data. 
When the admin is loged in, the website is the same but he gets the function to add/edit/delete the content of pages "Library" and "Cinematheque". These two pages content/data must be different and shouldn't interact with each other in any way.

Library: "Add" function lets you add either a new book or a new tag, tag can be things like "read" "2021" "project every house book" and so on. "Add a book" function makes you fill out a form listing Cover: ~blank text box~, Title: ~blank text box~, Author: ~blank text box~, Year: ~blank text box~, ISBN: ~blank text box~, AK's rating: ~blank text box~, Tags: ~blank text box~.

"Cover:" should be filled with the exact name of an image from your pc. AK's rating can be any number from 0 to 10. Example: "5", "6.5" and so on. The rest is self-explanatory. Tags can only be previously created ones in "add tags" section.

// FRONTEND - React //

- First page is "Home". About 2/10th of the upper screen is for showing page show/hide box (on the left) Look/desogn: simple white folder icon/emoji, page location (right after the page box buttom) and admin register box (on the right) BoxName: ADMIN.
- When opening the website you spawn in "Home" page and the page list is closed, when you click folder icon you see all the pages list in a "terminal directory folder tree" design like so:

├── Home/
│   ├── Portfolio/
│   │   ├── Education
│   │   ├── Experience
│   ├── Library
│   ├── Cinematheque

- (Page location example) folder_name/folder_name/file
- Clicking "ADMIN" gives you a new little square page for the admin to log in with only Title: AK LOGIN, Name: (blank text box), Name is AK and Password: (blank text box) it's a big password from an offline password mannager.
On every page after scrolling down, (down right of the screen) has a little box showing "up" sign meaning that you can jump back to top of the page.

Explanation of individual files:

- "Home" shows the "About" card explaining the website and why it exists.

- "Portfolio" shows one big card showing Contact: "my email" Github: "my github profile link" bellow is "Primary focus role:" foolowed by the desired job role in the moment and bellow "Also open to:" followed by less desired job roles I'm availabe to. Next to the role is a box saying "Resume", The box opens a new tab in visitor's browser with my resume pdf shared through drive in a viewer mode.

- "Education" shows a card with my university's logo picture on the left. Next to the picture is a title stating the university's name, On the right of the scree the time period is shown "2023 - Present" Bellow the title it says "Bachelor's Degree In Computer Science".

- "Experience" shows absolutelly all the experience I've had in life starting from the oldest going below to the newest like a life story of me.

- "Library" is AK's personal Goodreads. It shows search bar on top that has all the created tags listed. Below is just list of books with cover of the book and next to the cover is given the following information. Title: ~data~, Author: ~data~, ISBN: ~data~, Year: ~data~, AK's rating: ~data~/10, Tags: ~data~.


- "Cinematheque" Design and general function is the same as Library, with different tags and content. Each movie is a card with a cover, next to it is a list: Title: ~data~, Director: ~data~, Year: ~data~, AK's rating: ~data~, Tags: ~data~.







// AGAIN //
Admin Exists and she is the only one who can change the data though the website. For logging in the website needs only "Name" and "Password".
Website will have few pages but only two will be changeable through website/admin. Library and cinematheque. When the admin is logged in, she gets permission to add/delete content of both of these pages. 
Starting with "Library" page. Admin can only add/delete tags and books. adding tags is done by writing things like "2021" "to-read" "reading" in a blank text box in the wbsite and hitting create, you can't create two exact same tags or books. Adding books happens simply by writing down 