// App.jsx - EVERYTHING IN ONE FILE WITH BACKEND INTEGRATION
import React, { useState, useEffect, StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

// ============== CONFIGURATION ==============
const API_BASE = 'http://localhost:8080/api';

// ============== TAILWIND/CSS ==============
const styles = `
:root {
    --background: #121017;
    --primary: #b262fc;
    --foreground: #f3e8ff;
    --muted: #4c1d95;
    --muted-foreground: #a78bfa;
    --card-bg: rgba(0, 0, 0, 0.4);
    --border: rgba(172, 148, 244, 0.3);
}

* {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
}

body {
    background-color: var(--background);
    color: var(--foreground);
    font-family: 'Space Mono', monospace;
    line-height: 1.6;
    overflow-x: hidden;
}

.container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 40px 20px;
    position: relative;
    z-index: 10;
}

.scanline {
    background: linear-gradient(to bottom, transparent, transparent 50%, rgba(0,0,0,0.2) 50%, rgba(0,0,0,0.2));
    background-size: 100% 4px;
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    pointer-events: none;
    z-index: 50;
    opacity: 0.15;
}

.grid-bg {
    position: fixed;
    inset: 0;
    background-image: 
        linear-gradient(rgba(172, 148, 244, 0.1) 1px, transparent 1px),
        linear-gradient(90deg, rgba(172, 148, 244, 0.1) 1px, transparent 1px);
    background-size: 40px 40px;
    background-position: center bottom;
    opacity: 0.2;
    z-index: 0;
    pointer-events: none;
}

header {
    display: flex;
    gap: 40px;
    margin-bottom: 60px;
    padding-bottom: 40px;
    border-bottom: 2px solid var(--border);
    align-items: center;
}

@media (max-width: 768px) {
    header { flex-direction: column; text-align: center; }
}

.avatar-container {
    position: relative;
}

.avatar-frame {
    width: 160px;
    height: 160px;
    border: 4px solid var(--primary);
    padding: 4px;
    background: black;
    overflow: hidden;
    position: relative;
}

.avatar-frame img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    filter: contrast(125%) sepia(50%) hue-rotate(15deg);
}

.avatar-overlay {
    position: absolute;
    inset: 0;
    background: rgba(217, 70, 239, 0.1);
    mix-blend-mode: overlay;
}

.status-badge {
    position: absolute;
    bottom: -10px;
    right: -10px;
    background: black;
    border: 2px solid var(--primary);
    padding: 4px 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    font-weight: bold;
    box-shadow: 4px 4px 0 rgba(217, 70, 239, 0.5);
}

.status-dot {
    width: 8px;
    height: 8px;
    background: #22c55e;
    border-radius: 50%;
    animation: pulse 2s infinite;
}

@keyframes pulse {
    0% { transform: scale(1); opacity: 1; }
    50% { transform: scale(1.2); opacity: 0.7; }
    100% { transform: scale(1); opacity: 1; }
}

.text-glow {
    text-shadow: 0 0 10px var(--primary), 0 0 20px var(--primary);
}

h1 {
    font-family: 'VT323', monospace;
    font-size: 4rem;
    color: var(--primary);
    text-transform: uppercase;
    margin-bottom: 10px;
}

.role-text {
    font-family: 'VT323', monospace;
    font-size: 1.5rem;
    letter-spacing: 0.1em;
    opacity: 0.8;
}

.bio-text {
    margin-top: 20px;
    border-left: 4px solid var(--border);
    padding-left: 20px;
    font-style: italic;
    color: var(--muted-foreground);
}

.socials {
    margin-top: 30px;
    display: flex;
    gap: 15px;
    flex-wrap: wrap;
}

.socials a, .socials button {
    text-decoration: none;
    color: var(--primary);
    border: 1px solid var(--border);
    padding: 8px 16px;
    font-size: 12px;
    font-weight: bold;
    text-transform: uppercase;
    transition: all 0.3s;
    background: transparent;
    cursor: pointer;
    font-family: 'Space Mono', monospace;
}

.socials a:hover, .socials button:hover {
    background: var(--primary);
    color: black;
    box-shadow: 0 0 15px var(--primary);
}

main {
    display: grid;
    grid-template-columns: 250px 1fr;
    gap: 40px;
}

@media (max-width: 768px) {
    main { grid-template-columns: 1fr; }
}

nav {
    display: flex;
    flex-direction: column;
    gap: 15px;
}

nav button {
    background: transparent;
    border: 2px solid transparent;
    color: var(--muted-foreground);
    font-family: 'VT323', monospace;
    font-size: 1.5rem;
    padding: 15px 25px;
    text-align: left;
    cursor: pointer;
    transition: all 0.3s;
    text-transform: uppercase;
}

nav button.active {
    border-color: var(--primary);
    background: rgba(217, 70, 239, 0.1);
    color: var(--primary);
    box-shadow: 0 0 15px -5px var(--primary);
    transform: translateX(10px);
}

.retro-card {
    border: 2px solid var(--border);
    background: var(--card-bg);
    backdrop-filter: blur(4px);
    padding: 25px;
    position: relative;
    margin-bottom: 25px;
}

.retro-card h3 {
    font-family: 'VT323', monospace;
    font-size: 1.8rem;
    color: var(--primary);
    margin-bottom: 20px;
    border-bottom: 1px solid rgba(217, 70, 239, 0.2);
    padding-bottom: 10px;
    text-transform: uppercase;
}

.book-grid {
    display: grid;
    gap: 20px;
    grid-template-columns: 1fr 1fr;
}

@media (max-width: 600px) { 
    .book-grid { grid-template-columns: 1fr; } 
}

.book-card {
    border: 2px solid var(--border);
    background: var(--card-bg);
    padding: 15px;
    display: flex;
    gap: 15px;
    position: relative;
}

.book-cover {
    width: 80px;
    height: 120px;
    background: #1a1625;
    border: 1px solid var(--border);
    flex-shrink: 0;
    overflow: hidden;
}

.book-cover img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.book-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.book-title {
    font-family: 'VT323', monospace;
    font-size: 1.3rem;
    color: var(--primary);
}

.book-author {
    font-size: 0.9rem;
    opacity: 0.7;
}

.book-rating {
    color: #fbbf24;
    font-size: 0.9rem;
}

.book-tags {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
}

.tag {
    background: rgba(178, 98, 252, 0.2);
    border: 1px solid var(--primary);
    padding: 2px 8px;
    font-size: 0.75rem;
    color: var(--primary);
}

.admin-controls {
    position: absolute;
    top: 10px;
    right: 10px;
    display: flex;
    gap: 8px;
}

.icon-btn {
    background: transparent;
    border: 1px solid var(--border);
    color: var(--primary);
    padding: 4px 8px;
    cursor: pointer;
    font-size: 0.8rem;
    transition: all 0.3s;
}

.icon-btn:hover {
    background: var(--primary);
    color: black;
}

.add-btn {
    background: transparent;
    border: 2px solid var(--primary);
    color: var(--primary);
    padding: 12px 24px;
    font-family: 'VT323', monospace;
    font-size: 1.3rem;
    cursor: pointer;
    text-transform: uppercase;
    transition: all 0.3s;
    margin-bottom: 20px;
}

.add-btn:hover {
    background: rgba(178, 98, 252, 0.2);
    box-shadow: 0 0 15px rgba(178, 98, 252, 0.5);
}

.modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 20px;
}

.modal-content {
    background: var(--background);
    border: 2px solid var(--primary);
    padding: 30px;
    max-width: 500px;
    width: 100%;
    position: relative;
    max-height: 90vh;
    overflow-y: auto;
}

.modal-title {
    font-family: 'VT323', monospace;
    font-size: 2rem;
    color: var(--primary);
    margin-bottom: 20px;
    text-transform: uppercase;
}

.form-group {
    margin-bottom: 20px;
}

.form-label {
    display: block;
    margin-bottom: 8px;
    color: var(--muted-foreground);
    font-size: 0.9rem;
    text-transform: uppercase;
}

.form-input, .form-select {
    width: 100%;
    background: black;
    border: 1px solid var(--border);
    color: var(--foreground);
    padding: 10px;
    font-family: 'Space Mono', monospace;
    font-size: 0.9rem;
}

.form-input:focus, .form-select:focus {
    outline: none;
    border-color: var(--primary);
}

.form-actions {
    display: flex;
    gap: 15px;
    margin-top: 25px;
}

.btn-primary, .btn-secondary {
    flex: 1;
    padding: 12px 24px;
    font-family: 'VT323', monospace;
    font-size: 1.2rem;
    text-transform: uppercase;
    cursor: pointer;
    transition: all 0.3s;
    border: 2px solid var(--primary);
}

.btn-primary {
    background: var(--primary);
    color: black;
}

.btn-primary:hover {
    box-shadow: 0 0 20px var(--primary);
}

.btn-secondary {
    background: transparent;
    color: var(--primary);
}

.btn-secondary:hover {
    background: rgba(178, 98, 252, 0.2);
}

.error-message {
    color: #ef4444;
    font-size: 0.85rem;
    margin-top: 10px;
}

.loading {
    text-align: center;
    padding: 40px;
    font-family: 'VT323', monospace;
    font-size: 1.5rem;
    color: var(--primary);
}

footer {
    margin-top: 80px;
    padding-top: 40px;
    border-top: 1px solid var(--border);
    text-align: center;
    font-size: 12px;
    color: var(--muted-foreground);
}

.opacity-50 { opacity: 0.5; }
`;

// ============== MAIN APP COMPONENT ==============
function App() {
  // State
  const [activeTab, setActiveTab] = useState('about');
  const [isAdmin, setIsAdmin] = useState(false);
  const [token, setToken] = useState(localStorage.getItem('token') || '');
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [showBookModal, setShowBookModal] = useState(false);
  const [editingBook, setEditingBook] = useState(null);
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Check if already logged in on mount
  useEffect(() => {
    if (token) {
      setIsAdmin(true);
    }
  }, []);

  // Fetch books when library tab is active
  useEffect(() => {
    if (activeTab === 'library') {
      fetchBooks();
    }
  }, [activeTab]);

  // ============== API FUNCTIONS ==============
  const fetchBooks = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await fetch(`${API_BASE}/books`);
      if (response.ok) {
        const data = await response.json();
        setBooks(data || []);
      } else {
        setError('Failed to load books');
      }
    } catch (err) {
      setError('Error connecting to server');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleLogin = async (username, password) => {
    setError('');
    try {
      const response = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      if (response.ok) {
        const data = await response.json();
        setToken(data.token);
        localStorage.setItem('token', data.token);
        setIsAdmin(true);
        setShowLoginModal(false);
        return true;
      } else {
        setError('Invalid credentials');
        return false;
      }
    } catch (err) {
      setError('Error connecting to server');
      console.error(err);
      return false;
    }
  };

  const handleLogout = () => {
    setToken('');
    localStorage.removeItem('token');
    setIsAdmin(false);
  };

  const handleAddBook = async (bookData) => {
    setError('');
    try {
      const response = await fetch(`${API_BASE}/books`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(bookData)
      });

      if (response.ok) {
        await fetchBooks();
        setShowBookModal(false);
        return true;
      } else {
        setError('Failed to add book');
        return false;
      }
    } catch (err) {
      setError('Error connecting to server');
      console.error(err);
      return false;
    }
  };

  const handleUpdateBook = async (id, bookData) => {
    setError('');
    try {
      const response = await fetch(`${API_BASE}/books/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(bookData)
      });

      if (response.ok) {
        await fetchBooks();
        setShowBookModal(false);
        setEditingBook(null);
        return true;
      } else {
        setError('Failed to update book');
        return false;
      }
    } catch (err) {
      setError('Error connecting to server');
      console.error(err);
      return false;
    }
  };

  const handleDeleteBook = async (id) => {
    if (!confirm('Are you sure you want to delete this book?')) return;
    
    setError('');
    try {
      const response = await fetch(`${API_BASE}/books/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        await fetchBooks();
      } else {
        setError('Failed to delete book');
      }
    } catch (err) {
      setError('Error connecting to server');
      console.error(err);
    }
  };

  const tabs = ['about', 'portfolio', 'knowledge', 'library', 'films'];

  return (
    <>
      <style>{styles}</style>
      <div className="scanline"></div>
      <div className="grid-bg"></div>

      <div className="container">
        <header>
          <div className="avatar-container">
            <div className="avatar-frame">
              <img src="/AnaTheGreat.png" alt="Avatar" />
              <div className="avatar-overlay"></div>
            </div>
            <div className="status-badge">
              <div className="status-dot"></div>
              <span>ONLINE</span>
            </div>
          </div>

          <div className="profile-info">
            <h1 className="text-glow">AK_In_Web</h1>
            <p className="role-text">Owner // Georgia, Tbilisi</p>
            <p className="bio-text">I am AK. That's all I know. P.S. Please, check out ABOUT first to get context.</p>
            
            <div className="socials">
              {!isAdmin ? (
                <button onClick={() => setShowLoginModal(true)}>ADMIN</button>
              ) : (
                <button onClick={handleLogout}>LOGOUT</button>
              )}
            </div>
          </div>
        </header>

        <main>
          <nav>
            {tabs.map((tab, idx) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={activeTab === tab ? 'active' : ''}
              >
                {idx === 0 && '01. ABOUT'}
                {idx === 1 && '02. PORTFOLIO'}
                {idx === 2 && '03. KNOWLEDGE'}
                {idx === 3 && '04. LIBRARY'}
                {idx === 4 && '05. CINEMA'}
              </button>
            ))}
          </nav>

          <div className="content">
            {/* ABOUT SECTION */}
            {activeTab === 'about' && (
              <div className="retro-card">
                <h3>ABOUT THE WEBSITE</h3>
                <p>&gt; Website has 2 functions. AK's personal use and AK's portfolio. If you are an employer guest please check out PORTFOLIO. Everything's there for you.</p>
                <p style={{ marginTop: '15px' }}>&gt; This project is made by a talented, charming, full of perspective, young adult who believes in quote "If a field doesn't fascinate me, it's only because I don't know enough about it yet." -AK.</p>
                <p style={{ marginTop: '15px' }}>&gt; Despite the belief, human beings are limited creatures. Therefore, AK's main interests are only Mathematics, Tech (as you may have already guessed), Philosophy, Linguistics, History, Playing Guitar, Drawing, Origami.</p>
                <p style={{ marginTop: '15px' }}>&gt; As you can see, it's a lot. She's struggling, and hard times make a man do anything, including creating a website for managing the chaos and letting fellow humans know that exploring your interests doesn't have to be torturous.</p>
              </div>
            )}

            {/* PORTFOLIO SECTION */}
            {activeTab === 'portfolio' && (
              <div className="retro-card">
                <h3>PORTFOLIO</h3>
                <p>&gt; Coming soon...</p>
              </div>
            )}

            {/* KNOWLEDGE SECTION */}
            {activeTab === 'knowledge' && (
              <div className="retro-card">
                <h3>KNOWLEDGE</h3>
                <p>&gt; Coming soon...</p>
              </div>
            )}

            {/* LIBRARY SECTION */}
            {activeTab === 'library' && (
              <div>
                {isAdmin && (
                  <button className="add-btn" onClick={() => { setEditingBook(null); setShowBookModal(true); }}>
                    + Add New Book
                  </button>
                )}
                
                {error && <div className="error-message">{error}</div>}
                
                {loading ? (
                  <div className="loading">Loading books...</div>
                ) : books.length === 0 ? (
                  <div className="retro-card">
                    <p>&gt; No books yet. {isAdmin && 'Add your first book!'}</p>
                  </div>
                ) : (
                  <div className="book-grid">
                    {books.map(book => (
                      <div key={book.id} className="book-card">
                        <div className="book-cover">
                          {book.cover_image_url && <img src={book.cover_image_url} alt={book.title} />}
                        </div>
                        <div className="book-info">
                          <div className="book-title">{book.title}</div>
                          <div className="book-author">by {book.author}</div>
                          <div className="book-rating">{'★'.repeat(book.rating)}{'☆'.repeat(5 - book.rating)}</div>
                          {book.tags && book.tags.length > 0 && (
                            <div className="book-tags">
                              {book.tags.map((tag, i) => (
                                <span key={i} className="tag">{tag}</span>
                              ))}
                            </div>
                          )}
                        </div>
                        {isAdmin && (
                          <div className="admin-controls">
                            <button className="icon-btn" onClick={() => { setEditingBook(book); setShowBookModal(true); }}>✎</button>
                            <button className="icon-btn" onClick={() => handleDeleteBook(book.id)}>✕</button>
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}

            {/* FILMS SECTION */}
            {activeTab === 'films' && (
              <div className="retro-card">
                <h3>CINEMA</h3>
                <p>&gt; Coming soon...</p>
              </div>
            )}
          </div>
        </main>

        <footer>
          <p>EST. 2026 // All Rights Reserved</p>
          <p className="opacity-50">MADE BY AK</p>
        </footer>
      </div>

      {/* LOGIN MODAL */}
      {showLoginModal && (
        <LoginModal 
          onClose={() => setShowLoginModal(false)} 
          onLogin={handleLogin}
          error={error}
        />
      )}

      {/* BOOK MODAL */}
      {showBookModal && (
        <BookModal 
          book={editingBook}
          onClose={() => { setShowBookModal(false); setEditingBook(null); }}
          onSave={editingBook ? (data) => handleUpdateBook(editingBook.id, data) : handleAddBook}
          error={error}
        />
      )}
    </>
  );
}

// ============== LOGIN MODAL COMPONENT ==============
function LoginModal({ onClose, onLogin, error }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    await onLogin(username, password);
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <div className="modal-title">Admin Login</div>
        <div className="form-group">
          <label className="form-label">Username</label>
          <input 
            type="text" 
            className="form-input"
            value={username}
            onChange={e => setUsername(e.target.value)}
            placeholder="Enter username"
          />
        </div>
        <div className="form-group">
          <label className="form-label">Password</label>
          <input 
            type="password" 
            className="form-input"
            value={password}
            onChange={e => setPassword(e.target.value)}
            placeholder="Enter password"
          />
        </div>
        {error && <div className="error-message">{error}</div>}
        <div className="form-actions">
          <button className="btn-secondary" onClick={onClose}>Cancel</button>
          <button className="btn-primary" onClick={handleSubmit}>Login</button>
        </div>
      </div>
    </div>
  );
}

// ============== BOOK MODAL COMPONENT ==============
function BookModal({ book, onClose, onSave, error }) {
  const [formData, setFormData] = useState({
    title: book?.title || '',
    author: book?.author || '',
    cover_image_url: book?.cover_image_url || '',
    rating: book?.rating || 5,
    tags: book?.tags?.join(', ') || ''
  });

  const handleSubmit = async () => {
    const bookData = {
      ...formData,
      rating: parseInt(formData.rating),
      tags: formData.tags ? formData.tags.split(',').map(t => t.trim()).filter(Boolean) : []
    };
    await onSave(bookData);
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <div className="modal-title">{book ? 'Edit Book' : 'Add New Book'}</div>
        
        <div className="form-group">
          <label className="form-label">Title</label>
          <input 
            type="text" 
            className="form-input"
            value={formData.title}
            onChange={e => setFormData({...formData, title: e.target.value})}
            placeholder="Book title"
          />
        </div>

        <div className="form-group">
          <label className="form-label">Author</label>
          <input 
            type="text" 
            className="form-input"
            value={formData.author}
            onChange={e => setFormData({...formData, author: e.target.value})}
            placeholder="Author name"
          />
        </div>

        <div className="form-group">
          <label className="form-label">Cover Image URL</label>
          <input 
            type="text" 
            className="form-input"
            value={formData.cover_image_url}
            onChange={e => setFormData({...formData, cover_image_url: e.target.value})}
            placeholder="https://example.com/cover.jpg"
          />
        </div>

        <div className="form-group">
          <label className="form-label">Rating</label>
          <select 
            className="form-select"
            value={formData.rating}
            onChange={e => setFormData({...formData, rating: e.target.value})}
          >
            <option value="1">1 Star</option>
            <option value="2">2 Stars</option>
            <option value="3">3 Stars</option>
            <option value="4">4 Stars</option>
            <option value="5">5 Stars</option>
          </select>
        </div>

        <div className="form-group">
          <label className="form-label">Tags (comma separated)</label>
          <input 
            type="text" 
            className="form-input"
            value={formData.tags}
            onChange={e => setFormData({...formData, tags: e.target.value})}
            placeholder="Philosophy, Classic, Non-fiction"
          />
        </div>

        {error && <div className="error-message">{error}</div>}

        <div className="form-actions">
          <button className="btn-secondary" onClick={onClose}>Cancel</button>
          <button className="btn-primary" onClick={handleSubmit}>
            {book ? 'Update' : 'Add Book'}
          </button>
        </div>
      </div>
    </div>
  );
}

// ============== ENTRY POINT ==============
const rootElement = document.getElementById('root');
if (rootElement) {
  const root = createRoot(rootElement);
  root.render(
    <StrictMode>
      <App />
    </StrictMode>
  );
}

export default App;
