import { useState } from "react";
import { Link,NavLink,useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";

function Navbar({ theme, toggleTheme }) {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [query, setQuery] = useState("");

  const handleSearchSubmit = (e) => {
    e.preventDefault();

    if (query.trim()) {
      navigate(`/search?q=${encodeURIComponent(query.trim())}`);
    }
  };

  const handleGithubLogin = () => {
    window.location.href = "http://localhost:8080/auth/github";
  };

  return (
    <nav
      style={{
        display: "flex",
        alignItems: "center",
        padding: "10px 24px",
        borderBottom: "1px solid var(--color-border)",
        backgroundColor: "var(--color-bg-card)",
        gap: "24px",
      }}
    >
      {/* Logo */}
      <Link
        to="/"
        style={{
          fontWeight: "700",
          fontSize: "30px",
          textDecoration: "none",
          color: "var(--color-primary)",
          flexShrink: 0,
          marginRight: "150px",
        }}
      >
        Vouch
      </Link>

 {/* Navigation */}
<div
  style={{
    display: "flex",
    alignItems: "center",
    gap: "25px",
    flexShrink: 0,
  }}
>
  <NavLink
    to="/"
    style={({ isActive }) => ({
      textDecoration: isActive ? "underline" : "none",
      textUnderlineOffset: "4px",
      textDecorationThickness: "2px",
      color: isActive
        ? "#1d4ed8"
        : "var(--color-text-secondary)",
      fontSize: "14px",
      fontWeight: "600",
      transition: "color 0.15s ease",
    })}
    onMouseEnter={(e) => {
      e.currentTarget.style.color = "#1d4ed8";
    }}
    onMouseLeave={(e) => {
      e.currentTarget.style.color =
        e.currentTarget.getAttribute("aria-current")
          ? "#1d4ed8"
          : "var(--color-text-secondary)";
    }}
  >
    Home
  </NavLink>

  <NavLink
    to="/search"
    style={({ isActive }) => ({
      textDecoration: isActive ? "underline" : "none",
      textUnderlineOffset: "4px",
      textDecorationThickness: "2px",
      color: isActive
        ? "#1d4ed8"
        : "var(--color-text-secondary)",
      fontSize: "14px",
      fontWeight: "600",
      transition: "color 0.15s ease",
    })}
    onMouseEnter={(e) => {
      e.currentTarget.style.color = "#1d4ed8";
    }}
    onMouseLeave={(e) => {
      e.currentTarget.style.color =
        e.currentTarget.getAttribute("aria-current")
          ? "#1d4ed8"
          : "var(--color-text-secondary)";
    }}
  >
    Search
  </NavLink>

  {user && (
    <NavLink
      to={`/profile/${user.github_username}`}
      style={({ isActive }) => ({
        textDecoration: isActive ? "underline" : "none",
        textUnderlineOffset: "4px",
        textDecorationThickness: "2px",
        color: isActive
          ? "#1d4ed8"
          : "var(--color-text-secondary)",
        fontSize: "14px",
        fontWeight: "600",
        transition: "color 0.15s ease",
      })}
      onMouseEnter={(e) => {
        e.currentTarget.style.color = "#1d4ed8";
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.color =
          e.currentTarget.getAttribute("aria-current")
            ? "#1d4ed8"
            : "var(--color-text-secondary)";
      }}
    >
      Profile
    </NavLink>
  )}
</div>

      {/* Search */}
      <form
        onSubmit={handleSearchSubmit}
        style={{
          flex: 1,
          maxWidth: "360px",
          marginLeft: "auto",
        }}
      >
        <input
          type="text"
          placeholder="Search by name, project, or language..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="input-field"
          style={{
            padding: "8px 12px",
          }}
        />
      </form>

      {/* Right Side */}
      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: "10px",
          flexShrink: 0,
        }}
      >
        {/* Theme Toggle */}
        <button
          onClick={toggleTheme}
          aria-label="Toggle theme"
          title={theme === "light" ? "Dark mode" : "Light mode"}
          style={{
            width: "38px",
            height: "38px",
            borderRadius: "8px",
            border: "1px solid var(--color-border)",
            backgroundColor: "transparent",
            color: "var(--color-text-primary)",
            cursor: "pointer",
            fontSize: "16px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          {theme === "light" ? "☾" : "☀"}
        </button>

        {user ? (
          <>
            {/* New Post */}
            <button
              className="btn-primary"
              onClick={() => navigate("/create-post")}
            >
              New Post
            </button>

            {/* Profile Avatar */}
            <Link
              to={`/profile/${user.github_username}`}
              style={{
                display: "block",
              }}
            >
              <img
                src={user.avatar_url}
                alt={user.name}
                style={{
                  width: "40px",
                  height: "40px",
                  borderRadius: "50%",
                  objectFit: "cover",
                  display: "block",
                  border: "1px solid var(--color-border)",
                }}
              />
            </Link>

            {/* Logout */}
            <button
              onClick={logout}
              style={{
                backgroundColor: "transparent",
                color: "var(--color-text-secondary)",
                border: "1px solid var(--color-border)",
                padding: "8px 14px",
                borderRadius: "8px",
                fontWeight: "600",
                fontSize: "13px",
                cursor: "pointer",
              }}
            >
              Logout
            </button>
          </>
        ) : (
          <button className="btn-primary" onClick={handleGithubLogin}>
            Sign in
          </button>
        )}
      </div>
    </nav>
  );
}

export default Navbar;