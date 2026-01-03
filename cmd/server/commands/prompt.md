We want to add a job API. First create a job model/table with the follwomg properties:

- id (uuid) auto generated uuid
- title (string)
- description (text)
- location (string)
- salary (string) could be a range or a single value e.g $50000 or $50000 - $60000
- experience_level (string)
- skills (array of strings) 
- job_type (string) optional
- company (string)
- company_logo (string) optional
- created_at (timestamp)
- updated_at (timestamp)
- user (Relates to user model) one user can have many jobs (one to many relationship)

Create a job_route.go, job_model.go, job_repository.go, job_handler.go and job_service.go and add methods:
 - createJob that will allow authenticated user to create a job, the api request  should allow user to upload a company logo as part of the job creation request. The company logo should be uploaded to cloudinary and the url should be saved in the database as company_logo.
 - getAllJobs that will allow all authenticated users to get all jobs in the database.
 - getJobsByUser that will allow authenticated users to get all jobs created by them only.
 - getJobById that will allow all authenticated users to get a job by id.
 - updateJob that will allow authenticated user to update a job that they created. NB: Admin can update any job.
 - deleteJob that will allow authenticated user to delete a job that they created. NB: Admin can delete any job.

 Ensure proper error handling and validation for all methods.