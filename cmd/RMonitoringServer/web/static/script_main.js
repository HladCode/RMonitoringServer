if (!localStorage.getItem("API_BASE")) {
    const baseUrl = `${window.location.protocol}//${window.location.host}`;
    localStorage.setItem("API_BASE", baseUrl);
}

const API_BASE = localStorage.getItem("API_BASE");

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


// === DEVICES ===
async function getDevices() {
	const username = localStorage.getItem("username");
	if (!username) throw new Error("Имя пользователя не найдено.");

	const url = `${API_BASE}/user/getDevices`;
	const payload = {
		Username: username
	};

	let token = getAccessToken();

	try {
		let response = await fetch(url, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"Authorization": `Bearer ${token}`
			},
			body: JSON.stringify(payload)
		});

		if (response.status === 401) {
			token = await refreshJwtIfNeeded();
			return await getDevices(); // retry
		}

		if (!response.ok) {
			throw new Error(`Ошибка запроса: ${response.statusText}`);
		}

		const data = await response.json();
		localStorage.setItem("location_devices", JSON.stringify(data));
		return [...Object.keys(data), "-"];

	} catch (e) {
		throw new Error("Ошибка получения устройств: " + e.message);
	}
}

function getDevicesByLocation(location) {
	if (location === "-") return "";
	const raw = localStorage.getItem("location_devices");
	if (!raw) return "";
	const parsed = JSON.parse(raw);
	return parsed[location];
}

async function getSensors(device_id) {
	const url = `${API_BASE}/user/getSensors`;
	const payload = {
		Device_id: device_id
	};

	let token = getAccessToken();

	try {
		let response = await fetch(url, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"Authorization": `Bearer ${token}`
			},
			body: JSON.stringify(payload)
		});

		if (response.status === 401) {
			token = await refreshJwtIfNeeded();
			return await getSensors(device_id);
		}

		if (!response.ok) {
			throw new Error(`Ошибка получения сенсоров: ${response.statusText}`);
		}

		const data = await response.json();
		localStorage.setItem("sensors_meanings", JSON.stringify(data));
		return [...Object.keys(data), "-"];

	} catch (e) {
		throw new Error("Ошибка получения сенсоров: " + e.message);
	}
}

async function getSensorData(device_id, sensor_id, from, to) {
	const url = `${API_BASE}/user/getDataInInterval`;
	const payload = {
		ID: device_id,
		sensor_ID: sensor_id,
		from: from,
		to: to
	};

	let token = getAccessToken();

	try {
		let response = await fetch(url, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"Authorization": `Bearer ${token}`
			},
			body: JSON.stringify(payload)
		});

		if (response.status === 401) {
			token = await refreshJwtIfNeeded();
			return await getSensorData(device_id, sensor_id, from, to);
		}

		if (!response.ok) {
			throw new Error(`Ошибка получения данных: ${response.statusText}`);
		}

		return await response.json();

	} catch (e) {
		throw new Error("Ошибка получения данных сенсора: " + e.message);
	}
}

document.getElementById("loadLocations").addEventListener("click", async () => {
    try {
    const locations = await getDevices();
    const locSel = document.getElementById("locations");
    locSel.innerHTML = "";
    locations.forEach(loc => {
        const option = document.createElement("option");
        option.value = loc;
        option.textContent = loc;
        locSel.appendChild(option);
    });
    } catch (err) {
    alert("❌ Помилка завантаження локацій: " + err.message);
    }
});

document.getElementById("locations").addEventListener("change", () => {
    const location = document.getElementById("locations").value;
    const devices = getDevicesByLocation(location);
    const devSel = document.getElementById("devices");
    devSel.innerHTML = "";
    if (devices && devices.length) {
    devices.forEach(dev => {
        const option = document.createElement("option");
        option.value = dev;
        option.textContent = dev;
        devSel.appendChild(option);
    });
    }
});

document.getElementById("loadSensors").addEventListener("click", async () => {
    const device = document.getElementById("devices").value;
    if (!device) return alert("Виберіть пристрій!");
    try {
    const sensors = await getSensors(device);
    const senSel = document.getElementById("sensors");
    senSel.innerHTML = "";
    sensors.forEach(s => {
        const option = document.createElement("option");
        option.value = s;
        option.textContent = s;
        senSel.appendChild(option);
    });
    } catch (err) {
    alert("❌ Помилка завантаження сенсорів: " + err.message);
    }
});

// TODO: обробка кнопки loadChart
document.getElementById("loadChart").addEventListener("click", async () => {
    const device = document.getElementById("devices").value;
    const sensor = document.getElementById("sensors").value;
    const from = new Date(document.getElementById("start").value).toISOString();
    const to = new Date(document.getElementById("end").value).toISOString();

    if (!device || !sensor || !from || !to) {
    return alert("Виберіть пристрій, сенсор та вкажіть інтервал.");
    }

    try {
    const data = await getSensorData(device, sensor, from, to);
    drawChart(data);
    } catch (err) {
    alert("❌ Помилка побудови графіка: " + err.message);
    }
});

function drawChart(dataArray) {
    if (!dataArray || dataArray.length === 0) {
    return alert("Немає даних для графіка.");
    }

    google.charts.load('current', { packages: ['corechart'] });
    google.charts.setOnLoadCallback(() => {
    const data = new google.visualization.DataTable();
    data.addColumn('datetime', 'Час');
    data.addColumn('number', 'Значення');

    data.addRows(dataArray.map(item => [
        new Date(item.time),
        item.value
    ]));

    const chart = new google.visualization.LineChart(document.getElementById("chart_div"));
    chart.draw(data, {
        title: 'Графік сенсора',
        height: 400,
        width: '100%'
    });
    });
}