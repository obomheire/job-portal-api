lets made a major adjustment to our users and jobs models. We want the profile_picture and company_logo fieds to be a type {
    url: string
    public_id: string
}
This is necessary so that when we delete a user or job, we can also delete the profile_picture and company_logo from cloudinary using the public_id. Ensure to make a neccessary adjustment to all the affected API routes after this adjustment is made