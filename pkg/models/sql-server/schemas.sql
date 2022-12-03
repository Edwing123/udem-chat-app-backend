CREATE DATABASE [Nameless]
GO

CREATE TABLE [User] (
    [Id] INT IDENTITY(1, 1) PRIMARY KEY,

    [Name] NVARCHAR(40) NOT NULL,

    -- The length of the password is 60 characters
    -- long because the plain text password will be
    -- hashed using the hashing algorithm Bcrypt.
    [Password] CHAR(60) NOT NULL,

    [Birthdate] DATE NOT NULL,

    -- a UUID has a size of 36 characters.
    [Profile_Picture_Id] CHAR(36) NULL,

    -- Username must not be empty.
    CONSTRAINT [Check_User_Name_Not_Empty] CHECK (LEN(Name) > 0),

    -- Username must be unique.
    CONSTRAINT [Unique_User_Name] UNIQUE (Name),

    -- The of the profile picture must be unique.
    CONSTRAINT [Unique_User_Profile_Picture_Id] UNIQUE (Profile_Picture_Id)
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

    -- 300 is the maximum of characters per message.
    [Content] NVARCHAR(300) NOT NULL,

    [User_Id] INT NOT NULL,

    [Conversation_Id] INT NOT NULL,

    -- Foreign key references.
    CONSTRAINT [Foreign_Message_User_Id] FOREIGN KEY (User_Id) REFERENCES [User](Id),
    CONSTRAINT [Foreign_Message_Conversation_Id] FOREIGN KEY (Conversation_Id) REFERENCES [Conversation](Id),

    -- The message must not be empty.
    CONSTRAINT [Check_Message_Content_Not_Empty] CHECK (LEN(Content) > 0)
)
GO

CREATE TABLE [User_Join_Conversation] (
    [Id] INT IDENTITY(1, 1) PRIMARY KEY,

    [User_Id] INT NOT NULL,

    [Conversation_Id] INT NOT NULL,

    -- Foreign key references.
    CONSTRAINT [Foreign_User_Join_Conversation_User_Id] FOREIGN KEY (User_Id) REFERENCES [User](Id),
    CONSTRAINT [Foreign_User_Join_Conversation_Conversation_Id] FOREIGN KEY (Conversation_Id) REFERENCES [Conversation](Id)
)
