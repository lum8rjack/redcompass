// Utility function to update project status and unassign domains
export async function updateProjectStatus(pocketbase, project) {
  try {
    const response = await pocketbase.collection('Projects').update(project.id, {
      "Completed": !project.Completed
    });
    
    if(response) {
      // If marking as complete, unassign all domains and set Last_Used
      if (response.Completed) {
        const assignedDomains = await pocketbase.collection('Domains').getFullList({
          filter: `Assigned_Project = "${project.id}"`
        });

        // Unassign project from each domain and set Last_Used
        for (const domain of assignedDomains) {
            console.log(domain)
            await pocketbase.collection('Domains').update(domain.id, {
                "Assigned_Project": "",
                "Last_Used": project.id
            });
        }
      }
    }
    return response;
  } catch (error) {
    console.error('Error updating project status:', error);
    throw error;
  }
} 