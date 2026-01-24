const BOOK_TAGS = [
    "programming",
    "read in english",
    "book club",
    "project every house book",
    "mathematics",
    "philosophy",
    "podcast",
    "random",
    "want to read",
    "reading",
    "read"
];

const MOVIE_TAGS = [
    "🇹🇷",
    "🇰🇷",
    "drama",
    "horror",
    "psychological",
    "mystery",
    "action",
    "adventure",
    "fantasy",
    "comedy",
    "documentary",
    "animated",
    "cerebral",
    "visual-storytelling",
    "narrative",
    "artistic"
];

const BOOKS = [
    {
        id: 1,
        title: "The Design of Everyday Things",
        author: "Don Norman",
        isbn: "9780465050659",
        rating: 9.5,
        tags: validateTags(["design", "psychology"], BOOK_TAGS),
        cover: "design-of-everyday-things.webp",
        year: 1988,
        notes: "Essential guide to user-centered design principles."
    },
    {
        id: 2,
        title: "Refactoring",
        author: "Martin Fowler",
        rating: 9.0,
        year: 1999,
        tags: validateTags(["programming", "best-practices", "reference"], BOOK_TAGS),
        cover: "refactoring.webp",
        notes: "Techniques for restructuring code safely."
    }
];

const MOVIES = [
    {
        id: 1,
        title: "Shutter Island",
        director: "Martin Scorsese",
        year: 2010,
        rating: 8.5,
        tags: validateTags(["sci-fi", "cerebral", "artistic"], MOVIE_TAGS),
        poster: "shutterisland.webp",
        notes: "Mind-bending psychological thriller set in an asylum."
    },
    {
        id: 2,
        title: "Stalker",
        director: "Andrei Tarkovsky",
        rating: 10.0,
        year: 1979,
        tags: validateTags(["sci-fi", "cerebral", "artistic"], MOVIE_TAGS),
        poster: "stalker.webp",
        notes: "Contemplative journey into desire and meaning."
    }
];

function validateTags(itemTags, allowedTags) {
    return itemTags.filter(tag => allowedTags.includes(tag));
}

const WISDOM_ENTRIES = [
    {
        date: '24-Jan-2026',
        text: 'The best time to plant a tree was 20 years ago. The second best time is now.'
    },
    {
        date: '?-?-2026',
        text: 'Coming soon...'
    }
];

const state = {
    currentPage: 'about',
    libraryFilters: [],
    cinemaFilters: []
};

function getFirstLetter(str) {
    return str.charAt(0).toUpperCase();
}

function formatRating(rating) {
    const numRating = Number(rating);
    const formattedRating = numRating % 1 === 0 ? numRating.toFixed(0) : numRating.toFixed(1);
    return `AK's rating - ${formattedRating}/10`;
}

function getUsedTags(data, allTags) {
    const usedTags = new Set();
    data.forEach(item => {
        item.tags.forEach(tag => usedTags.add(tag));
    });
    
    return allTags.filter(tag => usedTags.has(tag));
}

function filterData(data, filters) {
    if (filters.length === 0) return data;
    return data.filter(item =>
        filters.every(filter => item.tags.includes(filter))
    );
}

function renderBookCard(book) {
    return `
        <div class="card">
            <img src="public/cover/${book.cover || 'default-cover.webp'}" 
                 alt="${book.title}" 
                 class="card-cover"
                 onerror="this.onerror=null; this.src='data:image/svg+xml;utf8,<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"120\" height=\"180\" viewBox=\"0 0 120 180\"><rect width=\"120\" height=\"180\" fill=\"%232a2a2a\"/><text x=\"60\" y=\"90\" text-anchor=\"middle\" fill=\"%23ff6b35\" font-family=\"VT323\" font-size=\"24\">${getFirstLetter(book.title)}</text></svg>';">
            <div class="card-content">
                <div class="card-title">${book.title}</div>
                <div class="card-meta">${book.author} • ${book.year} ${book.isbn ? `• ISBN: ${book.isbn}` : ''}</div>
                <div class="card-rating">${formatRating(book.rating)}</div>
                <div class="card-note">"${book.notes}"</div>
                <div class="card-tags">
                    ${book.tags.map(tag => `<span class="tag">#${tag}</span>`).join('')}
                </div>
            </div>
        </div>
    `;
}

function renderMovieCard(movie) {
    return `
        <div class="card">
            <img src="public/poster/${movie.poster || 'default-poster.webp'}" 
                 alt="${movie.title}" 
                 class="card-cover"
                 onerror="this.onerror=null; this.src='data:image/svg+xml;utf8,<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"120\" height=\"180\" viewBox=\"0 0 120 180\"><rect width=\"120\" height=\"180\" fill=\"%232a2a2a\"/><text x=\"60\" y=\"90\" text-anchor=\"middle\" fill=\"%2300d4ff\" font-family=\"VT323\" font-size=\"24\">${getFirstLetter(movie.title)}</text></svg>';">
            <div class="card-content">
                <div class="card-title">${movie.title}</div>
                <div class="card-meta">Directed by ${movie.director} • ${movie.year}</div>
                <div class="card-rating">${formatRating(movie.rating)}</div>
                <div class="card-note">"${movie.notes}"</div>
                <div class="card-tags">
                    ${movie.tags.map(tag => `<span class="tag">#${tag}</span>`).join('')}
                </div>
            </div>
        </div>
    `;
}

function renderLibrary() {
    const filteredData = filterData(BOOKS, state.libraryFilters);
    const usedTags = getUsedTags(BOOKS, BOOK_TAGS);
    const container = document.getElementById('library-cards');
    const badgesContainer = document.getElementById('filter-badges-library');
    const counterEl = document.getElementById('counter-library');
    const clearBtn = document.getElementById('clear-library');

    badgesContainer.innerHTML = BOOK_TAGS.map(tag => {
        const isUsed = usedTags.includes(tag);
        const isActive = state.libraryFilters.includes(tag);
        
        return `
            <button class="badge ${isActive ? 'active' : ''} ${!isUsed ? 'unused-tag' : ''}" 
                    data-tag="${tag}" 
                    data-type="library"
                    title="${!isUsed ? 'No books with this tag yet' : ''}"
                    ${!isUsed ? 'disabled' : ''}>
                #${tag}
            </button>
        `;
    }).join('');

    counterEl.textContent = `Showing ${filteredData.length} of ${BOOKS.length} volumes`;

    clearBtn.style.display = state.libraryFilters.length > 0 ? 'block' : 'none';

    if (filteredData.length === 0) {
        container.innerHTML = '<div class="no-matches">[ NO_MATCHES_FOUND ]</div>';
    } else {
        container.innerHTML = filteredData.map(book => renderBookCard(book)).join('');
    }
}

function renderCinema() {
    const filteredData = filterData(MOVIES, state.cinemaFilters);
    const usedTags = getUsedTags(MOVIES, MOVIE_TAGS);
    const container = document.getElementById('cinema-cards');
    const badgesContainer = document.getElementById('filter-badges-cinema');
    const counterEl = document.getElementById('counter-cinema');
    const clearBtn = document.getElementById('clear-cinema');

    badgesContainer.innerHTML = MOVIE_TAGS.map(tag => {
        const isUsed = usedTags.includes(tag);
        const isActive = state.cinemaFilters.includes(tag);
        
        return `
            <button class="badge ${isActive ? 'active' : ''} ${!isUsed ? 'unused-tag' : ''}" 
                    data-tag="${tag}" 
                    data-type="cinema"
                    title="${!isUsed ? 'No movies with this tag yet' : ''}"
                    ${!isUsed ? 'disabled' : ''}>
                #${tag}
            </button>
        `;
    }).join('');

    counterEl.textContent = `Showing ${filteredData.length} of ${MOVIES.length} volumes`;

    clearBtn.style.display = state.cinemaFilters.length > 0 ? 'block' : 'none';

    if (filteredData.length === 0) {
        container.innerHTML = '<div class="no-matches">[ NO_MATCHES_FOUND ]</div>';
    } else {
        container.innerHTML = filteredData.map(movie => renderMovieCard(movie)).join('');
    }
}

function renderWisdom() {
    const container = document.getElementById('wisdom-timeline');
    container.innerHTML = `
        <div class="timeline">
            ${WISDOM_ENTRIES.map(entry => `
                <div class="timeline-entry">
                    <div class="timeline-diamond"></div>
                    <div class="timeline-box">
                        <div class="timeline-date">${entry.date}</div>
                        <div class="timeline-text">${entry.text}</div>
                    </div>
                </div>
            `).join('')}
        </div>
    `;
}

function changePage(pageName) {
    document.querySelectorAll('.page').forEach(page => {
        page.classList.remove('active');
    });

    document.getElementById(pageName).classList.add('active');

    document.querySelectorAll('.nav-link').forEach(link => {
        link.classList.remove('active');
        if (link.dataset.page === pageName) {
            link.classList.add('active');
        }
    });

    state.currentPage = pageName;

    if (pageName === 'library') renderLibrary();
    if (pageName === 'cinema') renderCinema();
    if (pageName === 'wisdom') renderWisdom();
}

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            changePage(link.dataset.page);
        });
    });

    document.addEventListener('click', (e) => {
        if (e.target.classList.contains('badge') && !e.target.disabled) {
            const tag = e.target.dataset.tag;
            const type = e.target.dataset.type;
            const filterArray = type === 'library' ? state.libraryFilters : state.cinemaFilters;

            if (filterArray.includes(tag)) {
                filterArray.splice(filterArray.indexOf(tag), 1);
            } else {
                filterArray.push(tag);
            }

            if (type === 'library') renderLibrary();
            if (type === 'cinema') renderCinema();
        }
    });

    document.getElementById('clear-library').addEventListener('click', () => {
        state.libraryFilters = [];
        renderLibrary();
    });

    document.getElementById('clear-cinema').addEventListener('click', () => {
        state.cinemaFilters = [];
        renderCinema();
    });

    changePage('about');
});