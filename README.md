# Vouch

> **A developer portfolio and discovery platform built around your GitHub identity.**

**Vouch** helps developers turn their GitHub presence into a complete, discoverable developer profile. Developers can **sign in with GitHub, showcase their repositories, share posts, interact with other developers, and discover developers based on programming languages.**

## 💡Why Vouch?

A developer's work is often scattered across different platforms.

**GitHub** shows what they build, a **portfolio** describes their skills, and **social platforms** contain their thoughts and updates. Vouch brings these pieces together into one developer-focused platform.

With Vouch, developers can connect their **GitHub identity and repositories** to their profile, showcase their work, share their ideas through posts, and discover other developers based on their technical interests.

The idea is simple: **your actual work should be at the center of your developer identity.**

## Features

### GitHub Authentication

Developers can **sign in with GitHub** using OAuth. Vouch uses their GitHub identity to create or retrieve their account and connect their GitHub presence with their Vouch profile.

### Developer Profiles

Each developer gets a profile where they can **present their technical identity**, manage their information, and showcase their GitHub repositories.

### Repository Showcase

Vouch brings a developer's **public GitHub repositories** into their profile so visitors can explore the projects they have built and the technologies they work with.

### Developer Posts

Developers can share their **ideas, updates, projects, and experiences** through posts.

Posts can be:

- **Created**
- **Edited**
- **Deleted**

### Likes & Comments

Developers can **like and comment on posts**, making Vouch more interactive than a traditional portfolio.

### Developer Discovery

Vouch allows developers to **discover other developers based on programming languages** and explore their profiles, repositories, and posts.

### AI-Powered Features

Vouch integrates **Groq** to provide AI-powered functionality within the platform.

### Secure Backend

The backend uses **JWT authentication, protected API routes, and environment-based secret management** to secure application data and authenticated operations.

## Tech Stack

**Frontend:** React, JavaScript, HTML, CSS  
**Backend:** Go, Gin, GORM, REST APIs, JWT  
**Database:** PostgreSQL  
**Authentication:** GitHub OAuth, JWT  
**AI:** Groq  
**Deployment:** AWS EC2, Nginx

## Getting Started

### 📋Prerequisites

Make sure you have:

- **Go**
- **Node.js & npm**
- **PostgreSQL**
- **Git**
- **GitHub OAuth credentials**
- **Groq API key**

### Clone the Repository

    git clone https://github.com/bhumika019579/vouch.git
    cd vouch

### Backend

    cd server
    go mod download
    go run .

### Frontend

    cd client
    npm install
    npm run dev

### Environment Setup

Copy the example environment file and add your credentials:

    cp .env.example .env

**Never commit your `.env` file or expose your API keys.**

## 🌐Live Demo

**Coming soon**

## 👩‍💻Author

**Bhumika Chanchlani**

**GitHub:** https://github.com/bhumika019579  
**Portfolio:** https://bhumika019579.github.io/portfolio

