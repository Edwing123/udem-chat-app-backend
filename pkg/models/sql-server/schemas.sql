CREATE DATABASE [Nameless]
GO

CREATE TABLE [User] (
    [Id] INT IDENTITY(1, 1) PRIMARY KEY,

    -- The user name is unique because it will
    -- be used as a credential.
    [Name] NVARCHAR(40) NOT NULL UNIQUE,

    -- I guess this shouldn't be mentioned here.
    -- The length of the password is 60 bytes
    -- long because the plain text password will be
    -- hashed using the hashing algorithm Bcrypt.
    [Password] CHAR(60) NOT NULL,
    [BirthDate] DATE NOT NULL,

    -- a UUID is used for unique id values,
    -- and it has a size of 16 bytes.
    [Profile_Picture_Id] CHAR(16) NOT NULL UNIQUE,

    -- Username must not be empty.
    CONSTRAINT [Check_User_Name_Not_Empty] CHECK (LEN(Name) > 0)
)
GO

CREATE TABLE [Conversation] (
    [Id] INT IDENTITY(1, 1) PRIMARY KEY,
    [Created_At] DATETIME NOT NULL,

    -- The duration of the conversation in seconds.
    [Duration] SMALLINT NOT NULL,

    -- The duration of the conversation must not be negative.
    CONSTRAINT [Check_Conversation_Duration_Not_Negative] CHECK (Duration >= 0)
)
GO

CREATE TABLE [Message] (
    [Id] INT IDENTITY(1, 1) PRIMARY KEY,
    [Created_At] DATETIME NOT NULL,

    -- The maximum of characters per message.
    [Content] NVARCHAR(300) NOT NULL,
    [User_Id] INT NOT NULL,
    [Conversation_Id] INT NOT NULL,

    -- Foreigh key references.
    CONSTRAINT [Foreign_Message_User_Id] FOREIGN KEY [User](Id),
    CONSTRAINT [Foreign_Message_Conversation_Id] FOREIGN KEY [Conversation](Id),

    -- The message must not be empty.
    CONSTRAINT [Check_Message_Content_Not_Empty] CHECK (LEN(Content) > 0)
)
GO

CREATE TABLE [User_Join_Conversation] (
    [Id] INT IDENTITY(1, 1) PRIMARY KEY,
    [User_Id] INT NOT NULL,
    [Conversation_Id] INT NOT NULL,

    -- Foreigh key references.
    CONSTRAINT [Foreign_User_Join_Conversation_User_Id] FOREIGN KEY [User](Id),
    CONSTRAINT [Foreign_User_Join_Conversation_Conversation_Id] FOREIGN KEY [Conversation](Id),
)