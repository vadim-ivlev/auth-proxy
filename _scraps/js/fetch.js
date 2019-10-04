fetch('https://localhost:4400/apps/photoreports-admin/schema', {
        method: 'POST',
        credential: 'include',
        body: JSON.stringify({
          query: '',
          variables: {}
        })
      })
        .then(response => {
          if (response.ok) return response.json()
          new Error(response)
        })
        .then(res => { console.log(res) })
        .catch(error => { console.error(error) })