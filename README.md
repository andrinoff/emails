# Go SMTP Serverless Function for Vercel & iCloud+

[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8.svg?style=for-the-badge&logo=go)](https://go.dev/)
[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fandrinoff%2Femails)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

A simple, secure, and serverless Go function designed for deployment on Vercel. It acts as an API backend for a static website's contact form, sending emails through an iCloud+ custom domain email address.

## Purpose

This project provides a lightweight backend solution for static sites (like those built with Svelte, React, or plain HTML/JS) that need a "Get In Touch" or contact form. It avoids the need for a dedicated server or third-party email services by leveraging the SMTP server provided with an iCloud+ subscription.

## Features

- **Serverless**: Deploys as a single function on Vercel's free tier.
- **Secure**: Uses environment variables for credentials and an app-specific password, keeping secrets out of the code.
- **CORS Protected**: Includes middleware to only allow requests from your specified domains.
- **Lightweight**: Written in Go for fast cold starts and minimal resource usage.
- **Direct Reply**: Sets the `Reply-To` header, so you can reply directly to the person who contacted you.

---

## Setup and Deployment

Follow these steps to get your email API live.

### 1. Prerequisites

- An [iCloud+](https://www.apple.com/icloud/) subscription with a custom domain configured.
- A [Vercel](https://vercel.com) account.
- [Go](https://go.dev/dl/) installed on your machine.

### 2. Generate an iCloud App-Specific Password

You cannot use your regular Apple ID password. You must generate one for this application.

1.  Sign in to [appleid.apple.com](https://appleid.apple.com).
2.  Go to **Sign-In and Security** > **App-Specific Passwords**.
3.  Click **"Generate an app-specific password"**.
4.  Give it a label (e.g., `Vercel Contact Form`) and copy the generated password (`xxxx-xxxx-xxxx-xxxx`). **Save it somewhere safe.**

### 3. Clone and Configure the Project

Clone this repository to your local machine.

```bash
git clone https://github.com/andrinoff/emails.git
cd emails
```

The project is structured for Vercel. The function code is located at `/api/andrinoff/index.go`.

To use,

  a) change the name of the folder to whatever you want your endpoint to be located at.

  e.g.
    if folder is named `email` --> `https://<your-project-name>.vercel.app/api/email`
  b) Update CORS middleware (`line 25`) settings to your domain(s). (IMPORTANT!)
  c) Change email addresses at 79, 81 lines.


### 4. Set Environment Variables on Vercel

This is the most critical step. In your Vercel project dashboard, go to **Settings** -> **Environment Variables** and add the following:

| Key                           | Value                                                              | Description                                     |
| ----------------------------- | ------------------------------------------------------------------ | ----------------------------------------------- |
| `ICLOUD_AUTH_USER`            | Your primary Apple ID (e.g., `your-name@icloud.com`)               | The email used to log in to Apple's SMTP server. |
| `ICLOUD_APP_SPECIFIC_PASSWORD`| The password you generated in Step 2 (e.g., `xxxx-xxxx-xxxx-xxxx`) | The secure password for this application.       |

### 5. Deploy

Deploy the project to Vercel. If you have the Vercel CLI, you can run:

```bash
vercel --prod
```

Alternatively, push the repository to GitHub and link it in the Vercel dashboard. After setting the environment variables, Vercel will automatically build and deploy the function.

---

## API Usage

Once deployed, your function will be available at a URL like `https://<your-project-name>.vercel.app/api/<folder-name>`.

- **Method**: `POST`
- **Headers**:
  - `Content-Type: application/json`
- **Request Body** (JSON):

  ```json
  {
    "name": "John Doe",
    "email": "john.doe@example.com",
    "content": "Hello, I would like to get in touch!"
  }
  ```

### Example Frontend Fetch Request

Here is how you can call the API from your website's JavaScript:

```javascript
async function submitContactForm(name, email, content) {
  const endpoint = 'https://<your-project-name>.vercel.app/api/sendmail';

  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, email, content }),
    });

    const result = await response.json();

    if (!response.ok) {
      throw new Error(result.message || 'An error occurred.');
    }

    console.log('Success:', result.message);
    // Display success message to the user

  } catch (error) {
    console.error('Error:', error.message);
    // Display error message to the user
  }
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
