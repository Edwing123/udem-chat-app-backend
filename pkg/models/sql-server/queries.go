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
	INSERT INTO [User] ([Name], [Password], [Birthdate])
	VALUES(@Name, @Password, @Birthdate);
	`

	getUserById = `
	SELECT [Id], [Name], [Birthdate], [Profile_Picture_Id]
	FROM [User]
	WHERE [Id] = @Id;
	`

	getUserProfilePictureIdById = `
	SELECT [Profile_Picture_Id]
	FROM [User]
	WHERE [Id] = @Id;
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
