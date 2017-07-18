//Package render provides several types of renderable entities which are used throughout the code base
//
// In addition to entities the package also provides utilities to load images from files and load images
// from parts of files (subsprites such as is used in sprite sheets) as well as draw them.
//
// Renderable package has both simple entites and more complex entites that build off them.
// There is the Sprite that supports being loaded and drawn in a location, and  on the
// other end of the spectrum there is the Animation that can control how it draws sprite sheet sub components.
package render
