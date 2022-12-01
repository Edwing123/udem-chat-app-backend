package sqlserver

var (
	userId               = "Id"
	userName             = "Name"
	userPassword         = "Password"
	userBirthdate        = "Birthdate"
	userProfilePictureId = "Profile_Picture_Id"
)

const (
	insertUser = `
	INSERT INTO [User] ([Name], [Password], [Birthdate], [Profile_Picture_Id])
	VALUES(@Name, @Password, @Birthdate, @Profile_Picture_Id);
	`

	getUserById = `
	SELECT [Id], [Name], [Birthdate], [Profile_Picture_Id]
	FROM [User]
	WHERE [Id] = @Id;
	`

	getUserByName = `
	SELECT [Name]
	FROM [User]
	WHERE [Name] = @Name;
	`

	getUserPasswordById = `
	SELECT [Password]
	FROM [User]
	WHERE [Id] = @Id;
	`

	getUserIdAndPasswordByName = `
	SELECT [Id], [Password]
	FROM [User]
	WHERE [Name] = @Name;
	`

	updateUserPassword = `
	UPDATE [User]
	SET [Password] = @Password
	WHERE Id = @Id;
	`
)
