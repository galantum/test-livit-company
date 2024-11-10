const BASE_URL = 'http://127.0.0.1:8090/search';

export async function searchVideos(query, page = 1) {
  try {
    const response = await fetch(BASE_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ search: query, page }),
    });
    if (!response.ok) throw new Error('Failed to fetch videos');
    const data = await response.json();
    return data.data || [];
  } catch (error) {
    console.error('Error fetching videos:', error);
    return [];
  }
}
