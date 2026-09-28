function LoginPage() {
  const handleGithubLogin = () => {
    window.location.href = "http://localhost:8080/auth/github";
  };

  return (
    <div style={{ textAlign: "center", marginTop: "120px" }}>
      <h1 style={{ fontSize: "32px", fontWeight: "700" }}>Vouch</h1>
      <p style={{ color: "var(--color-text-secondary)", marginBottom: "24px" }}>
        A portfolio backed by proof, not claims.
      </p>
      <button className="btn-primary" onClick={handleGithubLogin}>
        Sign up with GitHub
      </button>
    </div>
  );
}

export default LoginPage;