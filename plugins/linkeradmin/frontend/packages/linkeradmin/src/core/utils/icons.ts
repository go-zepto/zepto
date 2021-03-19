export const loadMaterialIcons = () => {
  const link = document.createElement('link');
  link.setAttribute('rel', 'stylesheet');
  link.setAttribute('type', 'text/css');
  link.setAttribute('href', 'https://fonts.googleapis.com/icon?family=Material+Icons');
  document.head.appendChild(link);
};
