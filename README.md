# 🌦️ Weather CLI: Because Sometimes You’re *That* Desperate ⛅

**Imagine this**:  

You’re sitting in a room with no windows. The fluorescent lights hum softly, and you have no idea what’s happening in the outside world. Is it sunny with a chance of freedom? Or is a torrential rainstorm about to ruin the day you planned two weeks ago? Worse yet, your phone is dead. No apps. No browser. No desktop GUI. But wait… you’ve got **one thing**: a terminal and an irrational urge to know the current weather.

Enter: **Weather CLI**, the ultimate (and borderline unnecessary) tool for today's command-line enthusiasts. With Weather CLI, you’ll know if you need a raincoat, sunglasses, or sheer existential courage to step outside—all without ever leaving the comforting embrace of your terminal. 

---

## 🚀 Installation

You’ll need [Go](https://go.dev/doc/install) installed first. Once you’re set up, just run:

```bash
go install github.com/sss7526/weather@latest
```

Make sure your `$GOPATH/bin` (or `$GOBIN`) is in your `PATH` so you can summon the Weather CLI from anywhere. Add this to your `.bashrc`, `.zshrc`, or equivalent shell file if you haven’t already:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

That’s it. 🎉 You’re officially equipped with a weather-sniffing terminal wizard.

---

## 🖥️ Usage

Need to scratch that meteorological itch? Just run:

```bash
weather
```

Weather CLI will:
1. Detect your location using your public IP.
2. Fetch the most recent weather forecast using the **National Weather Service (NWS) API**.
3. Deliver it in a clean, no-nonsense format. Because the weather doesn’t have time for nonsense.

Example output:
```
Local Forecast:

This Afternoon 3 July 2025 at 01:00 PM EDT to 3 July 2025 at 06:00 PM EDT
Temp: 93°F | Wind: NE at 5 mph | Precip: 13%

Tonight 3 July 2025 at 06:00 PM EDT to 4 July 2025 at 06:00 AM EDT
Temp: 73°F | Wind: E at 1 to 5 mph | Precip: 15%

Independence Day 4 July 2025 at 06:00 AM EDT to 4 July 2025 at 06:00 PM EDT
Temp: 95°F | Wind: NE at 2 to 9 mph | Precip: 9%
...
```

No frills, no banner ads, no suspiciously upbeat weather anchors ✨.

---

## ⚙️ Features

- **Window-Free Weather**: Because sometimes your only “window” is `tmux`.
- **No Flood of Data**: Clean, concise weather information for each time period, because your terminal real estate is sacred.
- **Extremely Niche Appeal™**: Impress (or alienate) your coworkers by using the terminal to check if it’s raining.

---

## 👩‍💻 Contributing

Weather CLI is so lightweight that it practically disappears in the cloud. But if you’d like to enhance it (more features? International weather? Moon phase predictions?), feel free to fork it, submit a pull request, or open an issue. 

If you find bugs, it’s probably the API, not my code. Okay, okay, it might be my code. Let me know? Please?

---

## 💡 How It Works (Who Needs a Phone Anyway?)

- **Step 1**: Weather CLI asks the internet politely, “Where am I?”.
- **Step 2**: It leverages the U.S. **National Weather Service (NWS)** API for a detailed weather forecast. That’s right, a government-backed weather tool (you're welcome, taxpayers).
- **Step 3**: Results are displayed in a terminal-friendly, minimalist format. Every forecast is fetched just for YOU. Okay, and maybe everyone else with your public IP.

---

## 🤷 Why Did I Write This?

Because I could. And why not?  

- Phones die.  
- Windows aren't always available.  
- *But your terminal?* Your terminal never fails you.

Okay, maybe it wasn’t strictly necessary—kind of like a $600 mechanical keyboard to write emails—but now it’s here, and you’re clearly curious enough to install it. Plus, your terminal deserves a little joy once in a while.

Now go forth! Weather thyself.

☀️🌧️🌪️

--- 

## 📜 License

Licensed under the [MIT License](LICENSE). Use it, modify it, break it, or rewrite it in COBOL. Just don’t blame me if you read an outdated weather report and ruin your picnic. 

---

## Disclaimer

I know this is hard to believe but this cringe readme was written by a bot bc who's got time for that? Also if you're using a VPN it will pull the weather for the VPN location.

---