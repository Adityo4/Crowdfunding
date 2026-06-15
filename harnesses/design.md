# 🎨 Design Style Harness

> **Tujuan:** Memastikan setiap UI/UX yang dibangun memiliki visual language yang konsisten, estetis, dan premium.
> **Gunakan harness ini saat:** membangun antarmuka baru, membuat komponen UI, mendefinisikan design system, atau melakukan visual review.

---

## 🧭 Konteks untuk AI

```
Kamu adalah UI/UX designer dan front-end engineer yang berpengalaman.
Setiap keputusan visual harus:
1. Merujuk ke file referensi design system di folder references/ (jika tersedia)
2. Mengutamakan konsistensi — gunakan token yang sudah ada sebelum membuat baru
3. Mengikuti prinsip "Less is More" — kesederhanaan yang elegan lebih baik dari kompleksitas yang ramai
4. Selalu mempertimbangkan aksesibilitas (WCAG 2.1 AA minimum)
5. Dibangun mobile-first, lalu di-scale ke desktop

Sebelum membuat UI baru, SELALU tanya: "Sudah ada komponen atau token yang bisa dipakai?"
```

---

## 📂 Design References

Harness ini menggunakan **design reference files** yang tersimpan di folder:

```
harnesses/references/
└── <designname>-design.md    ← Spesifikasi visual spesifik per project/brand
```

### Cara Menggunakan References

Saat memulai sesi vibe coding dengan harness ini, sertakan file referensi yang relevan:

```
# Contoh penggunaan:
[Paste isi design.md ini] + [Paste isi references/genesis-design.md]
```

AI akan menggabungkan **prinsip umum** dari harness ini dengan **design token & visual identity spesifik** dari file referensi project-mu.

### Yang Harus Ada di File Reference

Setiap `<designname>-design.md` sebaiknya mendefinisikan:

```markdown
## Brand Identity
- Nama brand, tone of voice, kepribadian visual

## Color Tokens
- Primary, secondary, neutral, semantic colors (success, warning, error, info)

## Typography
- Font family, skala ukuran, weight yang dipakai

## Spacing & Layout
- Base unit, skala spacing, breakpoint

## Component Style
- Border radius, shadow, border style khas brand

## Inspirasi & Referensi
- Link atau deskripsi visual yang ingin dicapai
```

---

## 1. 🎨 Warna (Color)

### Prinsip
- Gunakan **design tokens** — jangan hardcode hex value langsung di komponen
- Setiap warna punya **semantic meaning** yang jelas
- Pastikan **contrast ratio** memenuhi WCAG: minimal 4.5:1 untuk teks normal, 3:1 untuk teks besar

### Struktur Token Warna

```css
:root {
  /* === BRAND COLORS === */
  --color-primary-50:  /* terlalu terang, untuk background */
  --color-primary-100:
  --color-primary-400:
  --color-primary-500: /* warna utama brand */
  --color-primary-600: /* hover state */
  --color-primary-900: /* terlalu gelap */

  /* === NEUTRAL === */
  --color-neutral-0:   /* putih */
  --color-neutral-100:
  --color-neutral-500: /* teks sekunder */
  --color-neutral-900: /* teks utama */
  --color-neutral-1000: /* hitam */

  /* === SEMANTIC === */
  --color-success: /* hijau */
  --color-warning: /* kuning/oranye */
  --color-error:   /* merah */
  --color-info:    /* biru */
}
```

### Checklist Warna
- [ ] Semua warna terdefinisi sebagai CSS custom property / design token
- [ ] Dark mode token tersedia jika diperlukan
- [ ] Contrast ratio diverifikasi untuk semua kombinasi teks-background
- [ ] Tidak ada warna yang di-hardcode di luar token

---

## 2. ✍️ Tipografi (Typography)

### Prinsip
- Maksimal **2 font family** dalam satu project (1 untuk heading, 1 untuk body)
- Gunakan **type scale** yang konsisten — jangan pakai ukuran sembarang
- **Line height** dan **letter spacing** sama pentingnya dengan ukuran font

### Type Scale yang Direkomendasikan

```css
:root {
  /* === FONT FAMILY === */
  --font-heading: 'Inter', sans-serif;   /* atau sesuai referensi */
  --font-body:    'Inter', sans-serif;

  /* === FONT SIZE (Type Scale) === */
  --text-xs:   0.75rem;   /*  12px */
  --text-sm:   0.875rem;  /*  14px */
  --text-base: 1rem;      /*  16px — base */
  --text-lg:   1.125rem;  /*  18px */
  --text-xl:   1.25rem;   /*  20px */
  --text-2xl:  1.5rem;    /*  24px */
  --text-3xl:  1.875rem;  /*  30px */
  --text-4xl:  2.25rem;   /*  36px */
  --text-5xl:  3rem;      /*  48px */

  /* === FONT WEIGHT === */
  --font-regular:   400;
  --font-medium:    500;
  --font-semibold:  600;
  --font-bold:      700;

  /* === LINE HEIGHT === */
  --leading-tight:  1.25;
  --leading-normal: 1.5;
  --leading-relaxed: 1.75;
}
```

### Checklist Tipografi
- [ ] Font di-load dari Google Fonts atau self-hosted (bukan system font default)
- [ ] Semua ukuran teks menggunakan token, bukan nilai arbitrary
- [ ] Heading hierarchy konsisten (h1 > h2 > h3)
- [ ] Panjang baris teks ideal: 45–75 karakter per baris

---

## 3. 📐 Spacing & Layout

### Prinsip
- Gunakan **base unit 4px** — semua spacing adalah kelipatan 4
- **Whitespace adalah desain** — jangan takut memberi ruang
- Layout harus **responsif** — mobile-first, bukan afterthought

### Spacing Scale

```css
:root {
  --space-1:  0.25rem;  /*  4px */
  --space-2:  0.5rem;   /*  8px */
  --space-3:  0.75rem;  /* 12px */
  --space-4:  1rem;     /* 16px */
  --space-5:  1.25rem;  /* 20px */
  --space-6:  1.5rem;   /* 24px */
  --space-8:  2rem;     /* 32px */
  --space-10: 2.5rem;   /* 40px */
  --space-12: 3rem;     /* 48px */
  --space-16: 4rem;     /* 64px */
  --space-20: 5rem;     /* 80px */
  --space-24: 6rem;     /* 96px */
}
```

### Breakpoint

```css
/* Mobile-first breakpoints */
--bp-sm:  640px;   /* Tablet kecil */
--bp-md:  768px;   /* Tablet */
--bp-lg:  1024px;  /* Desktop kecil */
--bp-xl:  1280px;  /* Desktop */
--bp-2xl: 1536px;  /* Desktop besar */
```

### Checklist Layout
- [ ] Grid system terdefinisi (12 kolom atau CSS Grid custom)
- [ ] Max-width konten terdefinisi (biasanya 1200px–1440px)
- [ ] Semua spacing menggunakan token, bukan nilai arbitrary
- [ ] Layout berfungsi di mobile (320px) hingga wide desktop (1920px)

---

## 4. 🧱 Komponen UI

### Prinsip
- **Atomic Design** — bangun dari atom (button, input) ke molekul (form field) ke organisme (form)
- Setiap komponen harus punya **semua state**: default, hover, focus, active, disabled, error
- Komponen harus **accessible** — keyboard navigable, ARIA label yang tepat

### Komponen Dasar yang Harus Ada

| Komponen | States yang Diperlukan |
|---|---|
| Button | Default, Hover, Active, Disabled, Loading |
| Input | Default, Focus, Error, Disabled, Read-only |
| Select / Dropdown | Default, Open, Selected, Disabled |
| Checkbox / Radio | Unchecked, Checked, Indeterminate, Disabled |
| Modal / Dialog | Open, Close animation |
| Toast / Alert | Success, Warning, Error, Info |
| Card | Default, Hover (jika clickable) |
| Badge / Tag | Berbagai variasi warna |

### Checklist Komponen
- [ ] Semua state visual terdefinisi dan diimplementasi
- [ ] Focus state visible dan jelas (outline atau ring)
- [ ] Komponen bisa digunakan dengan keyboard saja
- [ ] ARIA attributes tepat (`aria-label`, `role`, `aria-expanded`, dll)
- [ ] Komponen tidak bergantung pada warna saja untuk menyampaikan informasi

---

## 5. ✨ Motion & Animasi

### Prinsip
- Animasi harus **purposeful** — memperjelas UX, bukan sekedar dekoratif
- Durasi **singkat dan konsisten** — UI bukan film
- Hormati preferensi `prefers-reduced-motion`

### Durasi yang Direkomendasikan

```css
:root {
  --duration-instant:  50ms;   /* Feedback micro (highlight, ripple) */
  --duration-fast:    150ms;   /* Hover, tooltip appear */
  --duration-normal:  250ms;   /* Panel open, dropdown */
  --duration-slow:    400ms;   /* Modal, page transition */
  --duration-slower:  600ms;   /* Complex animation, onboarding */

  --ease-default:     cubic-bezier(0.4, 0, 0.2, 1);  /* Material standard */
  --ease-in:          cubic-bezier(0.4, 0, 1, 1);
  --ease-out:         cubic-bezier(0, 0, 0.2, 1);
  --ease-spring:      cubic-bezier(0.34, 1.56, 0.64, 1); /* Bouncy */
}
```

### Checklist Animasi
- [ ] Semua transisi menggunakan token durasi dan easing
- [ ] `prefers-reduced-motion` direspek — animasi di-disable atau dikurangi
- [ ] Tidak ada animasi yang berjalan terus-menerus tanpa tujuan
- [ ] Loading state menggunakan skeleton screen, bukan hanya spinner

---

## 6. 🌑 Dark Mode

### Strategi Implementasi

```css
/* Gunakan CSS custom properties + media query atau class */

@media (prefers-color-scheme: dark) {
  :root {
    --color-background: #0f0f0f;
    --color-surface:    #1a1a1a;
    --color-text:       #f5f5f5;
    /* override token yang perlu berubah */
  }
}

/* Atau dengan class untuk toggle manual: */
[data-theme="dark"] {
  --color-background: #0f0f0f;
  /* ... */
}
```

### Checklist Dark Mode
- [ ] Semua warna menggunakan token (sehingga mudah di-override)
- [ ] Tidak ada warna yang hardcode — semuanya lewat token
- [ ] Gambar dan ilustrasi masih terlihat di dark background
- [ ] Shadow disesuaikan (shadow gelap kurang efektif di dark mode, gunakan glow/border)

---

## 7. ♿ Aksesibilitas (A11y)

### Prinsip
- Aksesibilitas bukan tambahan — ini bagian dari desain yang baik
- Target minimum: **WCAG 2.1 Level AA**

### Checklist Aksesibilitas
- [ ] Contrast ratio teks: ≥ 4.5:1 (normal), ≥ 3:1 (besar/bold)
- [ ] Semua elemen interaktif bisa diakses via keyboard (Tab, Enter, Space, Arrow)
- [ ] Focus indicator terlihat jelas di semua elemen interaktif
- [ ] Gambar punya `alt` text yang deskriptif (atau `alt=""` jika dekoratif)
- [ ] Form memiliki label yang terhubung ke input (`for` / `aria-labelledby`)
- [ ] Error message jelas dan spesifik (bukan hanya warna merah)
- [ ] Heading hierarchy logis (tidak skip dari h1 ke h4)
- [ ] ARIA digunakan dengan benar — jangan tambah ARIA jika HTML semantik sudah cukup

---

## 8. 📋 Design Review Checklist (Pre-Launch)

### Visual Consistency
- [ ] Semua warna menggunakan design token
- [ ] Semua spacing menggunakan skala yang terdefinisi
- [ ] Tipografi konsisten dengan type scale
- [ ] Tidak ada nilai arbitrary yang scattered di kode

### Komponen & Interaksi
- [ ] Semua state komponen sudah diimplementasi (hover, focus, disabled, error)
- [ ] Animasi menggunakan token durasi dan easing
- [ ] Loading dan empty state terdefinisi untuk semua area konten

### Responsivitas
- [ ] Tampilan diuji di breakpoint: 320px, 768px, 1024px, 1440px
- [ ] Tidak ada horizontal scroll di mobile
- [ ] Touch target minimal 44×44px di mobile

### Aksesibilitas
- [ ] Contrast ratio lulus WCAG AA
- [ ] Navigasi keyboard berfungsi di seluruh halaman
- [ ] Screen reader test dilakukan (VoiceOver / NVDA)

### Reference Check
- [ ] File referensi di `references/<designname>-design.md` sudah diperbarui
- [ ] Token baru ditambahkan ke file referensi

---

## 📚 Referensi

- [WCAG 2.1 Guidelines](https://www.w3.org/TR/WCAG21/)
- [Refactoring UI](https://www.refactoringui.com/)
- [Design Tokens W3C](https://design-tokens.github.io/community-group/format-spec/)
- [Inclusive Components](https://inclusive-components.design/)
- [Contrast Checker](https://webaim.org/resources/contrastchecker/)

---

*Design Style Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
*Untuk token & visual identity spesifik project, lihat [references/](./references/)*
