import {tags} from '@params';

const KeyBackspace = 8;
const KeyEnter     = 13;
const KeyEsc       = 27;
const KeyLeft      = 37;
const KeyUp        = 38;
const KeyRight     = 39;
const KeyDown      = 40;
const KeyR         = 82;
const KeyS         = 83;

const paginator = {
  next: () => {
    let a = document.querySelectorAll('[aria-label="Next"]')[0];
    if (!a) {
      return;
    }

    window.location = a.href;
  },
  prev: () => {
    let a = document.querySelectorAll('[aria-label="Previous"]')[0];
    if (!a) {
      return;
    }

    window.location = a.href;
  }
};

let search = {
  enabled: false,
  accuracy: 0.8,
  val: "",
  el: undefined,

  init: function() {
    // Insert search into the page, below the title
    this.el = document.createElement("div");
    this.el.setAttribute("id", "search");
    let title = document.getElementsByClassName("title")[0];
    if (title) {
      title.parentElement.appendChild(this.el);
    }
  },
  enable: function() {
    this.enabled = true;
    this.el.classList.remove("hidden");
    this.setResults();
  },
  disable: function() {
    this.enabled = false;
    this.val = "";
    this.el.innerHTML = "";
    this.el.classList.add("hidden");
  },
  append: function(s) {
    this.val += s.toLowerCase();
    this.setResults();
  },
  del: function() {
    this.val = this.val.slice(0, -1);
    this.setResults();
  },
  setResults: function() {
    let html = "";
    tags.filter((tag) => this.fuzzy(tag, this.val, this.accuracy)).forEach((tag) => {
      // Wrap each matching tag letter in a span
      let displayTag = "";
      for (let t=0; t < tag.length; t++) {
        let tagLetter = tag[t],
            match = false;
        for (let v=0; v < this.val.length; v++) {
          let queryLetter = this.val[v];
          if (tagLetter === queryLetter) {
            match = true;
          }
        }
        if (match) {
          displayTag += '<span class="highlight">' + tagLetter + '</span>';
        } else {
          displayTag += tagLetter;
        }
      }
      html += '<a href="/tags/' + tag + '"><div>' + displayTag + '</div></a>'
    });
    this.el.innerHTML = html;
  },
  bestMatch: function() {
    let matches = tags.filter((tag) => this.fuzzy(tag, this.val, this.accuracy));
    return matches.length ? matches[0] : undefined;
  },
  fuzzy: function(string, term, ratio) {
    string = string.toLowerCase();
    let compare = term.toLowerCase();
    let matches = 0;
    if (string.indexOf(compare) > -1) return true; // covers basic partial matches
    for (let i = 0; i < compare.length; i++) {
      string.indexOf(compare[i]) > -1 ? matches += 1 : matches -=1;
    }
    return (matches/string.length >= ratio || term == "")
  }
}

// Global keybindings
document.addEventListener("keydown", (e) => {
  switch(e.which) {
    // Navigation controls
    case KeyLeft:
      paginator.prev();
      e.preventDefault();
      break;

    case KeyRight:
      paginator.next();
      e.preventDefault();
      break;

    case KeyR: // r: random 
      if (!search.enabled) {
        window.location = "/random";
      } else {
        search.append(String.fromCharCode(e.keyCode));
      }
      break;

    // Search controls
    case KeyS:
      if (!search.enabled) {
        search.enable();
        console.log("searchEnabled: ", search.enabled);
        return;
      }
      search.append(String.fromCharCode(e.keyCode));
      break;

    case KeyEsc:
      search.disable();
      console.log("searchEnabled: ", search.enabled);
      break;

    case KeyBackspace:
      e.preventDefault();
      if (search.enabled) {
        search.del();
      }
      break;

    case KeyEnter:
      if (!search.enabled) {
        return;
      }

      let match = search.bestMatch();
      if (match) {
        window.location = "/tags/" + match;
      }
      break;

    default:
      if (!search.enabled) {
        return;
      }
      search.append(String.fromCharCode(e.keyCode));

  }
}, false);

// Finally, inject search
search.init();
