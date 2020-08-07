const get = async <T>(url: string): Promise<{ response: T; error: any }> => {
  return fetch(url, {
    method: "GET",
  })
    .then(async (response) => {
      const data = await response.json();
      if (!response.ok) {
        const error = (data && data.message) || response.status;
        return Promise.reject({ response: null, error });
      }

      return { response: data, error: null };
    })
    .catch((error) => {
      return { response: null, error };
    });
};

const apiClient = {
  get,
};

export default apiClient;
