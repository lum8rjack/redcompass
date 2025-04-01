import PocketBase from 'pocketbase';

// Initialize PocketBase instance
const url = window.location.origin;
const pocketbase = new PocketBase(url);

export default pocketbase;