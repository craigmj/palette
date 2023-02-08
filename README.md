# palette

Palette is a command line tool to generate semantic css structures for sites and webcomponents.

## Understanding the problem

I'm building increasingly complicated sites in CSS, DOM and Webcomponents, and even without webcomponents, CSS too quickly becomes a setting of particular values (colors for example), on a per-selector basis. For example, in the header links:

```css
header a {
	color: red;
	&:hover {
		color: green;
	}
}
```

I'm not a designer, but a programmer, and will change my mind about design quite a few times before a designer looks at the work and changes the palette, the font sizes, and so on.

So I started to think about a 'semantic css', using CSS variables:

```css
header a {
	color: var(--link-color);
	&:hover {
		color: var(--link-color-hover);
	}
}
```

I can define all my colours in one place, and can change the colour values as necessary. When I'm using webcomponents and shadow-DOM, this becomes particulary useful, because there isn't another way (aside for `::part`) to style a webcomponent. What's more, because of the 'cascading' nature of CSS variables, it allows me to define general semantic ideas (the color of a link), and redefine these in different areas of page (for instance, the navigation bar, the header, or the main work area), just by redefining the CSS variables in those areas:

```css
header {
	--link-color: red;
	--link-color-hover: green;
}
main {
	--link-color: blue;
	--link-color-hover: green;
}
```

What's more, if I use those CSS variables in my webcomponents, they will style themselves appropriately to their environment, without being infected by other styling of the environment.

## An experimental solution

I start off with the variables I want to set:

  color
  background-color
  font-size

and sometimes these might be qualified by interactive conditions, such as:

  color-hover
  background-color-hover

Each of these is a CSS variable that I will set for my components.

(No, this is too complicated - what I really need to do is set this with a SCSS macro,
but even if I do so, I can use my variable system... but they why not _just_ use my variable system, and keep all the values in SCSS: because that makes it _impossible_ to override with themes...: what we are actually defining here is something closer to a CSS theming system... a basic one, since we can't override everything, but we can do basic overrides that will allow some cosmetic changes to be incorporated with CSS variables.)

I define my variables and give each a default value.

variables:
  color: black;
  background-color: white;
  font-size: 14pt;


Then there are aspects. I'm thinking of something like 'aspect-orientated programming', except for styling.

One aspect of styling is 'context'. This is _where_ the variable is applied. For instance, there is a 'link' or 'text' or a 'list-item', or a 'table-row'. These can be defined by the developer. There should be as few as possible, but as many as necessary, and their distinction should be clear.

Another aspect might be 'role'. For instance, a button might be the 'primary' button or the 'secondary' button. It might be a 'warning' or a 'success' notification.

With every aspect, the rule that seems best is to have as few as possible, as distinct as spossible.

