// Utility function to update project status and unassign domains
export async function updateProjectStatus(pocketbase, project) {
  try {
    const response = await pocketbase.collection('Projects').update(project.id, {
      "Completed": !project.Completed
    });
    
    return response;
  } catch (error) {
    console.error('Error updating project status:', error);
    throw error;
  }
} 