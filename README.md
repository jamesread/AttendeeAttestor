# AttendeeAttestor

**Private proof of presence.**

AttendeeAttestor is a free and open-source, privacy-preserving QR ticket verification system designed for real-world events such as discos, festivals, conferences, and private venues.

It allows staff to verify that an attendee is valid **without ever learning who that person is**, and continues working even when the internet goes down.

---

## ✅ Key Features

- 🔐 **Cryptographically signed QR codes**
- 📵 **Fully offline verification**
- 🕵️ **No personal data required**
- 🔁 **Replay protection (anti double-entry)**
- 🧑‍✈️ **Designed for real door staff, not engineers**
- 🧩 **Open-source and self-hostable**

---

## 🧠 How It Works (High Level)

1. Tickets are issued as **signed QR codes**.
2. Each QR code contains:
   - Event ID
   - Random ticket ID
   - Valid time window
   - Ticket type (e.g. General, VIP)
3. No names, emails, phone numbers, or accounts are ever embedded.
4. Staff scan with the AttendeeAttestor app.
5. The app verifies the **signature offline** using a public key.
6. The app confirms:
   - The QR is authentic
   - The event matches
   - The ticket is within its valid time window
   - The ticket has not already been used

✅ Green = Admit  
❌ Red = Reject  

No central database lookup is required for verification.

---

## 🔒 Privacy by Design

AttendeeAttestor intentionally avoids:
- Tracking
- Profiles
- Accounts for attendees
- Facial recognition
- Network dependency

The scanner proves *validity*, not *identity*.

---

## 🔁 Preventing Re-Use

Each scanner maintains a local cache of used ticket IDs.  
When connectivity is available, scanners can optionally sync used entries to prevent multi-door reuse.

---

## 🛠️ Project Structure (Planned)

- `attendeeattestor-issuer` – Ticket generation & signing
- `attendeeattestor-scanner` – Staff scanning app (offline-first)
- `attendeeattestor-core` – Shared verification logic
- `docs/` – Protocol & threat model

---

## 🪪 Use Cases

- Nightclubs & discos
- Conferences & trade shows
- Festivals
- Private events
- Community meetups

---

## 📜 License

AGPL .

---

## 🚦 Project Status

Early design & protocol definition.

If you're interested in:
- contributing,
- reviewing the crypto workflow,
- or building a scanner UI,
