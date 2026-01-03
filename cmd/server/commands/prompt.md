Now, let's implement the following methods:
1. GetAllUsers (return all the users in the database). Only admin can access this method
2. DeleteUser (delete a user by id from the database). This should be done using a transaction. Add a deleteAsset method in the pkg/cloudinary/cloudinary.go file that will delete a file/asset from cloudinary using the public_id. Use it to delete the user's profile picture from cloudinary. Then delete the user from the database. Also ensure that all the jobs created by the user are cascaded on delete.
3. ForgotPassword. User make request  to reset password with email address. Check that the email exist, if yes, generate a 6 digit random number and return to the user. Later we will update this so that the 6 digit number is sent to the user's email instead. Add a second method call ResetPassword that will accept the users email, new password and the 6 digit number. Validate the number and if it is correct and not expired, update the user's password. Throw error if the number is incorrect or expired. NB: the generated 6 digit number shiuld expire after 60 minutes. Ensure that, plain password is not saved to the database.
4. ChangePassword by user. Allow login user to change their password. User should provide their current password and new password.
5. ChangeUserPassword by admin. Allow admin to change a user's password. User should provide the user's id and new password. 

GetAllUsers and DeleteUser should be implemented in the user's, route, handlers, service and repository files.
ForgotPassword, ChangeUserPassword, ChangePassword and ResetPassword should be implemented in the auth's, route, handlers, service and repository files.
