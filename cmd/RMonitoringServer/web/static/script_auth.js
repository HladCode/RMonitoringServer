if (!localStorage.getItem("API_BASE")) {
    const baseUrl = `${window.location.protocol}//${window.location.host}`;
    localStorage.setItem("API_BASE", baseUrl);
}

const API_BASE = localStorage.getItem("API_BASE");

// === CONFIG ===
// function saveApiBaseUrl(url) {
// 	localStorage.setItem("API_BASE", url);
// }

// function API_BASE {
// 	return localStorage.getItem("API_BASE");
// }

// === TOKENS ===
function saveTokens(tokens) {
	localStorage.setItem("access_token", tokens.access_token);
	localStorage.setItem("refresh_token", tokens.refresh_token);
	localStorage.setItem("username", tokens.username);
}

function clearTokens() {
	localStorage.removeItem("access_token");
	localStorage.removeItem("refresh_token");
	localStorage.removeItem("username");
}

function getAccessToken() {
	return localStorage.getItem("access_token");
}

// === AUTH ===
async function registerUser(login, email, password) {
	const url = `${API_BASE}/auth/register`;
	const payload = {
		Login: login,
		email: email,
		password: password
	};

	try {
		const response = await fetch(url, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(payload)
		});
		return await response.text();
	} catch (e) {
		throw new Error("Ошибка при регистрации: " + e.message);
	}
}

async function loginUser(login, password) {
	const url = `${API_BASE}/auth/login`;
	const payload = {
		Login: login,
		password: password
	};

	try {
		const response = await fetch(url, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(payload)
		});

		if (!response.ok) {
			const errorText = await response.text();
			throw new Error(`Ошибка входа: ${response.status} ${errorText}`);
		}

		const text = await response.text();
        const cleanedText = text.replace(/200\s*$/, "");
        const result = JSON.parse(cleanedText);


		if (!result.access_token || !result.refresh_token) {
			throw new Error("Сервер не вернул токены.");
		}

		saveTokens({
			access_token: result.access_token,
			refresh_token: result.refresh_token,
			username: login
		});

		//return result.access_token;

	} catch (e) {
		throw new Error("Ошибка авторизации: " + e.message);
	}
}

async function refreshJwtIfNeeded() {
	const refresh_token = localStorage.getItem("refresh_token");
	const username = localStorage.getItem("username");

	if (!refresh_token || !username) {
		throw new Error("Нет сохранённых данных авторизации.");
	}

	const url = `${API_BASE}/auth/refresh`;
	const payload = {
		Login: username,
		Token: refresh_token
	};

	try {
		const response = await fetch(url, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(payload)
		});

		if (!response.ok) {
			clearTokens();
			throw new Error("Ошибка обновления токена. Авторизуйтесь заново.");
		}

		const text = await response.text();
        const cleanedText = text.replace(/200\s*$/, "");
        const result = JSON.parse(cleanedText);

		if (!result.access_token) {
			throw new Error("Сервер не вернул новый access_token.");
		}

		localStorage.setItem("access_token", result.access_token);
		return result.access_token;

	} catch (e) {
		clearTokens();
		throw new Error("Ошибка при обновлении JWT: " + e.message);
	}
}

document.getElementById("loginBtn").addEventListener("click", async () => {
    const name = document.getElementById("loginName").value;
    const pass = document.getElementById("loginPass").value;
    if (!name || !pass) return alert("Заповніть логін та пароль");
    try {
    await loginUser(name, pass);
    alert("✅ Успішний вхід!");
    } catch (err) {
    alert("❌ Помилка авторизації: " + err.message);
    }
});

document.getElementById("registerBtn").addEventListener("click", async () => {
    const name = document.getElementById("regName").value;
    const email = document.getElementById("regEmail").value;
    const pass = document.getElementById("regPass").value;
    const confirm = document.getElementById("regConfirm").value;

    if (!name || !email || !pass || !confirm) return alert("Заповніть всі поля");
    if (pass !== confirm) return alert("Паролі не співпадають!");

    try {
    await registerUser(name, email, pass);
    alert("✅ Успішна реєстрація!");
    } catch (err) {
    alert("❌ Помилка реєстрації: " + err.message);
    }
});