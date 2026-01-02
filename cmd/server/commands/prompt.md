Now, let's add a UpdateUserProfile method to update user's profile. 
NB:
1. The method will accept a user ID param
2. Only a user with isAdmin true can update any user's profile, a user that have isAdmin set to false can only update their own profile
3. Admin can update any properties of the user object  except password (username, email, isAdmin, profilePicture) but a user that is not an admin can only update 3 fields (username, email, profilePicture) from thier profile
4. The UpdateUserProfile should be athenticated